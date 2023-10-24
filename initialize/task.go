package initialize

import (
	"context"
	"data-transform/global"
	"data-transform/model/po"
	"data-transform/model/vo"
	"encoding/json"
	"fmt"
	"github.com/go-redis/redis/v8"
	rabbitmq "github.com/kgip/go-rabbit-template"
	"github.com/robfig/cron/v3"
	"strconv"
	"time"
)

func Task() {
	c := cron.New()
	//重发超时消息
	c.AddFunc("* */3 * * * ?", func() {
		//查询过期未确认的消息
		var list = make([]*po.MqMessage, 0)
		if err := global.DB.Model(po.MqMessageModel).Where("service_name = ? and state = ? and updated < ?", "data-transform", 1, time.Now().Sub(time.Unix(int64(time.Minute/time.Second), 0))).Find(&list).Error; err != nil {
			global.LOG.Error(err.Error())
		} else {
			for i := 0; i < len(list); i++ {
				importKgVo := &vo.ImportKgVo{}
				if err = json.Unmarshal([]byte(list[i].Content), importKgVo); err != nil {
					global.LOG.Error(err.Error())
					continue
				}
				if err = global.MQ.SimplePublish("uploadTask.exchange", "uploadTask.builder", importKgVo, &rabbitmq.CorrelationData{ID: list[i].MessageId}); err != nil {
					global.LOG.Error(err.Error())
				} else {
					list[i].Retry++
					if err = global.DB.Updates(list[i]).Error; err != nil {
						global.LOG.Error(err.Error())
					}
				}
			}
		}
	})
	//检查中断的任务
	c.AddFunc("* */2 * * * ?", func() {
		var list = make([]*po.Task, 0)
		if err := global.DB.Model(po.TaskModel).Where("state in ? and status = ?", []int{1, 2}, 1).Find(&list).Error; err != nil {
			global.LOG.Error(err.Error())
		} else {
			for i := 0; i < len(list); i++ {
				//检查是否有节点正在执行该任务
				if res, err := global.Redis.Eval(context.Background(), "", []string{}).Result(); err != nil {
					global.LOG.Error(err.Error())
				} else {
					if res == 0 { //没有有节点正在执行该任务
						if maxIndex, err := global.Redis.Get(context.Background(), fmt.Sprintf("task:max_index:{%d}", list[i].ID)).Result(); err != nil {
							global.LOG.Error(err.Error())
						} else {
							if indexes, err := global.Redis.SMembers(context.Background(), fmt.Sprintf("task:uploaded_file_indexes:{%d}", list[i].ID)).Result(); err != nil {
								global.LOG.Error(err.Error())
							} else {
								startIndex, _ := strconv.ParseInt(maxIndex, 10, 64)
								var uploadedFileIndexes = map[string]struct{}{}
								for _, index := range indexes {
									n, _ := strconv.ParseInt(index, 10, 64)
									if n > startIndex {
										uploadedFileIndexes[index] = struct{}{}
									}
								}
								TaskService.DoUpload(list[i], 1, list[i].StartIndex, int(startIndex), list[i].LastIndex, uploadedFileIndexes)
							}
						}
					}
				}
			}
		}
	})
	//定时同步数据交易系统状态
	c.AddFunc("*/20 * * * * ?", func() {
		//查询所有处于导入中的任务
		var list = make([]*po.Task, 0)
		if err := global.DB.Model(po.TaskModel).Where("state = ? and status = ?", 3, 1).Find(&list).Error; err != nil {
			global.LOG.Error(err.Error())
		} else {
			for i := 0; i < len(list); i++ {
				if taskState, err := DataExchangeService.SyncTaskState(list[i].Token, list[i].ProdId); err != nil {
					global.LOG.Error(err.Error())
				} else {
					//更新状态
					if taskState.State == 3 && taskState.Status == 1 {
						_, err = global.Redis.Pipelined(context.Background(), func(p redis.Pipeliner) error {
							if result, err := p.HGet(context.Background(), fmt.Sprintf("task:progress:{%d}", list[i].ID), "count").Result(); err != nil {
								return err
							} else if count, _ := strconv.ParseInt(result, 10, 64); int(count) < taskState.Count {
								p.HMSet(context.Background(), fmt.Sprintf("task:progress:{%d}", list[i].ID), "count", taskState.Count, "progress", taskState.Progress)
							}
							return nil
						})
						if err != nil {
							global.LOG.Error(err.Error())
						}
					} else if taskState.Status == 2 { //执行失败
						err = global.DB.Updates(&po.Task{ID: list[i].ID, Status: 2, ErrorDetail: taskState.Cause}).Error
						if err != nil {
							global.LOG.Error(err.Error())
						}
					} else { //任务完成
						err = global.DB.Updates(&po.Task{ID: list[i].ID, Status: 3, State: 4}).Error
						if err != nil {
							global.LOG.Error(err.Error())
						}
					}
				}
			}
		}
	})
}
