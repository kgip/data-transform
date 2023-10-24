package service

import (
	"context"
	"crypto/md5"
	"data-transform/global"
	"data-transform/http"
	"data-transform/model/po"
	"data-transform/model/vo"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-redis/redis/v8"
	rabbitmq "github.com/kgip/go-rabbit-template"
	"github.com/kgip/redis-lock/lock"
	"github.com/robfig/cron/v3"
	uuid "github.com/satori/go.uuid"
	"os"
	"strconv"
	"sync"
	"time"
)

var ServerId string

func init() {
	hostname, _ := os.Hostname()
	ServerId = hostname + ":" + uuid.NewV4().String()
}

type TaskService struct {
	DataExchangeService http.DataExchangeService
	BuilderService      http.BuilderService
	OssGateway          http.OssGateway
	C                   *cron.Cron
}

func (service *TaskService) Upload(kgId string) {
	//查询上传绑定
	var bindingModel = &po.KgUploadBinding{}
	err := global.DB.Model(po.KgUploadBindingModel).Select("token_info_id", "max_index").Where("kg_id = ?", kgId).First(bindingModel).Error
	if err != nil {
		panic(err)
	}
	if bindingModel.TokenInfoId <= 0 {
		panic(errors.New("not bind production"))
	}
	//查询token信息
	var tokenInfoModel = &po.TokenInfo{}
	err = global.DB.Model(po.TokenInfoModel).Select("token", "prod_id").Where("id = ?", bindingModel.TokenInfoId).First(tokenInfoModel).Error
	if err != nil {
		panic(err)
	}
	//校验产品是否存在
	state, err := service.DataExchangeService.CheckProduction(tokenInfoModel.Token, tokenInfoModel.ProdId)
	if err != nil {
		panic(err)
	} else if state != 1 {
		panic(errors.New(fmt.Sprintf("production %s state error", tokenInfoModel.ProdId)))
	}
	err = global.LockOperator.Lock("lock:transform:addKgUploadTask:"+kgId, lock.Context())
	if err != nil {
		panic(err)
	}
	defer global.LockOperator.Unlock("lock:transform:addKgUploadTask:" + kgId)
	//校验图谱状态
	kgInfo, err := service.BuilderService.GetKgInfo(kgId)
	if err != nil {
		panic(err)
	}
	if kgInfo.State != 2 || kgInfo.MaxIndex <= bindingModel.MaxIndex {
		panic(errors.New("kg state error"))
	}
	//查询是否有未完成的任务
	var count int64
	err = global.DB.Model(po.TaskModel).Count(&count).Where("kg_id = ? and prod_id = ? and state in ?", kgId, tokenInfoModel.ProdId, []int{1, 2}).Error
	if err != nil {
		panic(err)
	}
	if count > 0 {
		panic(errors.New("kg is uploading"))
	}
	//创建任务
	task := &po.Task{
		KgId:       kgId,
		ProdId:     tokenInfoModel.ProdId,
		Token:      tokenInfoModel.Token,
		State:      1,
		Status:     1,
		StartIndex: bindingModel.MaxIndex + 1,
		LastIndex:  kgInfo.MaxIndex,
	}
	tx := global.DB.Create(task)
	if tx.Error != nil {
		panic(err)
	}
	if tx.RowsAffected <= 0 {
		panic(errors.New("create task error"))
	}
	//通知交易平台开始上传图谱
	err = service.DataExchangeService.Notice(tokenInfoModel.Token, tokenInfoModel.ProdId, 1)
	if err != nil {
		panic(errors.New("notice data exchange error: " + err.Error()))
	}
	//异步任务开始上传图谱
	go service.DoUpload(task, kgInfo.IsMetaUpdate, bindingModel.MaxIndex+1, bindingModel.MaxIndex+1, kgInfo.MaxIndex, nil)
}

func (service *TaskService) DoUpload(task *po.Task, isMetaUpdate, initIndex, startIndex, endIndex int, uploadedFileIndexes map[string]struct{}) {
	//心跳任务
	idleCtx, idleCancel := context.WithCancel(context.Background())
	entryId, _ := service.C.AddFunc("* */1 * * * ?", func() {
		for i := 0; i < 3; i++ {
			_, err := global.Redis.Pipelined(context.Background(), func(p redis.Pipeliner) error {
				key := fmt.Sprintf("idle:upload_task:%d", task.ID)
				if res, err := p.Get(context.Background(), key).Result(); err != nil {
					return err
				} else if res == "" {
					if err = p.Set(context.Background(), key, ServerId, 3*time.Minute).Err(); err != nil {
						return err
					}
				} else if res == ServerId {
					if err = p.Expire(context.Background(), key, 3*time.Minute).Err(); err != nil {
						return err
					}
				} else { //其它节点已经在执行该任务，取消当前节点任务执行
					idleCancel()
				}
				return nil
			})
			if err == nil {
				break
			} else {
				global.LOG.Error(err.Error())
			}
		}
	})
	defer service.C.Remove(entryId)
	if task.State == 1 {
		//检查图谱元数据是否更新，如果更新，需要先同步元数据
		if isMetaUpdate == 2 {
			if ok, err := service.SyncKgMetaData(task); !ok {
				if err != nil {
					global.LOG.Error(err.Error())
				}
				return
			}
		}
		//更新任务状态
		if err := global.DB.Updates(&po.Task{ID: task.ID, State: 2}).Error; err != nil {
			global.LOG.Error(err.Error())
			err = global.DB.Updates(&po.Task{ID: task.ID, Status: 2, ErrorDetail: "update task error: " + err.Error()}).Error
			if err != nil {
				global.LOG.Error(err.Error())
			}
			return
		}
	}
	//并发上传ds
	taskChannel := make(chan int, global.Config.Task.Concurrency*2)
	ctx, cancel := context.WithCancel(context.Background())
	errCtx, errCancel := context.WithCancelCause(context.Background())
	wg := &sync.WaitGroup{}
	for i := 0; i < global.Config.Task.Concurrency; i++ {
		go func(ctx context.Context) {
			for {
				select {
				case index := <-taskChannel:
					filename := fmt.Sprintf("%s_%d", task.ProdId, index)
					//下载文件
					url, err := service.OssGateway.GetDownloadInfo(filename)
					if err != nil {
						global.LOG.Error(err.Error())
						errCancel(err)
						wg.Done()
						return
					}
					file, hash, err := service.OssGateway.DownloadFile(url)
					if err != nil {
						global.LOG.Error(err.Error())
						errCancel(err)
						wg.Done()
						return
					}
					//校验文件hash
					h := md5.New()
					h.Write(file)
					hashRes := h.Sum(nil)
					if string(hashRes) != hash {
						global.LOG.Error(fmt.Sprintf("file %s hash not match", filename))
						errCancel(err)
						wg.Done()
						return
					}
					//上传文件
					err = service.DataExchangeService.UploadFile(file, task.Token, task.ProdId, filename, hash)
					if err != nil {
						global.LOG.Error(err.Error())
						errCancel(err)
						wg.Done()
						return
					}
					//上传成功，修改上传进度
					_, err = global.Redis.Pipelined(context.Background(), func(p redis.Pipeliner) error {
						if err = p.SAdd(context.Background(), fmt.Sprintf("task:uploaded_file_indexes:{%d}", task.ID), index).Err(); err != nil {
							return err
						}
						res, err := p.HIncrBy(context.Background(), fmt.Sprintf("task:progress:{%d}", task.ID), "count", 1).Result()
						if err != nil {
							return err
						}
						if err = p.HSet(context.Background(), fmt.Sprintf("task:progress:{%d}", task.ID), "progress", int(res)*100/(endIndex-startIndex+1)).Err(); err != nil {
							return err
						}
						return nil
					})
					if err != nil {
						global.LOG.Error(err.Error())
					}
					wg.Done()
				default:
					select {
					case <-ctx.Done():
						return
					default:
						time.Sleep(time.Second)
					}
				}
			}
		}(ctx)
	}
	for i := startIndex; i <= endIndex; i += global.Config.Task.Concurrency {
		select {
		case <-errCtx.Done():
			err := global.DB.Updates(&po.Task{ID: task.ID, Status: 2, ErrorDetail: "upload file error: " + errCtx.Err().Error()}).Error
			if err != nil {
				global.LOG.Error(err.Error())
				return
			}
		case <-idleCtx.Done():
			return
		default:
			var j int
			for j = i; j <= endIndex && j < i+global.Config.Task.Concurrency; j++ {
				if _, ok := uploadedFileIndexes[fmt.Sprintf("%d", j)]; !ok {
					wg.Add(1)
					taskChannel <- j
				}
			}
			wg.Wait()
			//修改进度
			_, err := global.Redis.Pipelined(context.Background(), func(p redis.Pipeliner) error {
				if err := p.Set(context.Background(), fmt.Sprintf("task:max_index:{%d}", task.ID), j-1, 48*time.Hour).Err(); err != nil {
					return err
				}
				if result, err := p.HGet(context.Background(), fmt.Sprintf("task:progress:{%d}", task.ID), "count").Result(); err != nil {
					return err
				} else if count, _ := strconv.ParseInt(result, 10, 64); int(count) < j-initIndex {
					p.HMSet(context.Background(), fmt.Sprintf("task:progress:{%d}", task.ID), "count", j-initIndex, "progress", (j-initIndex)*100/(endIndex-initIndex+1))
				}
				if err := p.Unlink(context.Background(), fmt.Sprintf("task:uploaded_file_indexes:{%d}", task.ID)).Err(); err != nil {
					return err
				}
				return nil
			})
			if err != nil {
				global.LOG.Error(err.Error())
			}
		}
	}
	cancel()
	//上传完成
	//通知数据交易系统开始导入图谱
	err := service.DataExchangeService.ImportKg(task.Token, task.ProdId, startIndex, endIndex)
	if err != nil {
		err = global.DB.Updates(&po.Task{ID: task.ID, Status: 2, ErrorDetail: "update task error: " + err.Error()}).Error
		if err != nil {
			global.LOG.Error(err.Error())
		}
		return
	}
	//更新状态
	err = global.DB.Updates(&po.Task{ID: task.ID, State: 3}).Error
	if err != nil {
		err = global.DB.Updates(&po.Task{ID: task.ID, Status: 2, ErrorDetail: "update task error: " + err.Error()}).Error
		if err != nil {
			global.LOG.Error(err.Error())
		}
		return
	}
}

func (service *TaskService) SyncKgMetaData(task *po.Task) (bool, error) {
	data, err := service.BuilderService.ExportKgMetaData(task.KgId)
	if err != nil {
		global.LOG.Error(err.Error())
		//更新状态
		err = global.DB.Updates(&po.Task{ID: task.ID, Status: 2, ErrorDetail: "export kg metadata error: " + err.Error()}).Error
		return false, err
	}
	//上传元数据
	if err = service.DataExchangeService.SyncKgMetaData(data); err != nil {
		global.LOG.Error(err.Error())
		err = global.DB.Updates(&po.Task{ID: task.ID, Status: 2, ErrorDetail: "export kg metadata error: " + err.Error()}).Error
		return false, err
	}
	return true, nil
}

func (service *TaskService) ImportKg(importKgVo *vo.ImportKgVo) {
	//创建任务
	tx := global.DB.Begin()
	if err := tx.Create(&po.Task{KgId: importKgVo.ProdId, ProdId: importKgVo.ProdId, Status: 1, State: 3, StartIndex: importKgVo.StartIndex, LastIndex: importKgVo.EndIndex}).Error; err != nil {
		tx.Rollback()
		panic(err)
	}
	bytes, err := json.Marshal(importKgVo)
	if err != nil {
		tx.Rollback()
		panic(err)
	}
	messageId := uuid.NewV4().String()
	message := &po.MqMessage{ServiceName: "data-transform", MessageId: messageId, Content: string(bytes)}
	//发送消息通知builder开始导入
	if err = tx.Create(message).Error; err != nil {
		tx.Rollback()
		panic(err)
	}
	tx.Commit()
	if err = global.MQ.SimplePublish("uploadTask.exchange", "uploadTask.builder", importKgVo, &rabbitmq.CorrelationData{ID: messageId}); err != nil {
		global.LOG.Error(err.Error())
	}
}
