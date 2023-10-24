package initialize

import (
	"data-transform/global"
	"data-transform/model/po"
	rabbitmq "github.com/kgip/go-rabbit-template"
)

func MQ() {
	global.MQ, _ = rabbitmq.NewRabbitTemplate("amqp://guest:guest@192.168.206.129:5672/", rabbitmq.Config{EnablePublisherConfirm: true, EnablePublisherReturns: true})

	if err := global.MQ.ExchangeDeclare("uploadTask.exchange", rabbitmq.ExchangeTopic, true, false, false, true, nil); err != nil {
		panic(err)
	}
	if err := global.MQ.QueueDeclare("uploadTask.queue", true, false, false, true, nil); err != nil {
		panic(err)
	}
	if err := global.MQ.QueueBind("uploadTask.queue", "uploadTask.builder", "uploadTask.exchange", true, nil); err != nil {
		panic(err)
	}

	global.MQ.RegisterConfirmCallback(func(ack bool, DeliveryTag uint64, correlationData *rabbitmq.CorrelationData) {
		if err := global.DB.Model(po.MqMessageModel).Where("service_name = ? and message_id = ?", "data-transform", correlationData.ID).Update("state", 2).Error; err != nil {
			global.LOG.Error(err.Error())
		}
	})
}
