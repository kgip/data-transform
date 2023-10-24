package initialize

import (
	"data-transform/global"
	"data-transform/model/po"
	"data-transform/model/vo"
	"encoding/json"
	rabbitmq "github.com/kgip/go-rabbit-template"
	"github.com/robfig/cron/v3"
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
}
