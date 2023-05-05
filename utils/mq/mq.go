package mq

import (
	"errors"
	"log"
	"sync"
	"time"
)

var (
	QueueNotExist = errors.New("message queue does not exist")
)

const (
	defaultQueueSize = 100
	defaultTopicName = "default"
)

type MqOperator interface {
	SendMessage(topic string, msg interface{}, expire time.Duration) (err error)
	RegistryMessageHandler(topics []string, handler func(msg interface{})) (error error)
}

type LocalMQ struct {
	lock   *sync.Mutex
	queues map[string]chan interface{}
}

func NewLocalMQ() *LocalMQ {
	return &LocalMQ{lock: &sync.Mutex{}, queues: make(map[string]chan interface{})}
}

func (mq *LocalMQ) tryGetQueue(topic string) chan interface{} {
	if queue, ok := mq.queues[topic]; ok && queue != nil {
		return queue
	}
	return nil
}

func (mq *LocalMQ) AddQueue(topic string, size int64) *LocalMQ {
	if queue := mq.tryGetQueue(topic); queue == nil {
		mq.lock.Lock()
		defer mq.lock.Unlock()
		if queue = mq.tryGetQueue(topic); queue == nil {
			queue = make(chan interface{}, size)
			mq.queues[topic] = queue
		}
	}
	return mq
}

func (mq *LocalMQ) AddDefaultSizeQueue(topic string) *LocalMQ {
	return mq.AddQueue(topic, defaultQueueSize)
}

func (mq *LocalMQ) handleError(topic string, msg interface{}, handler func(msg interface{})) {
	defer func() {
		if e := recover(); e != nil {
			time.Sleep(3 * time.Second)
			mq.queues[topic] <- msg //处理失败时重回队列
			log.Println(e)
		}
	}()
	handler(msg)
}

func (mq *LocalMQ) SendMessage(topic string, msg interface{}, expire time.Duration) (err error) {
	if queue, ok := mq.queues[topic]; ok && queue != nil {
		if expire <= 0 {
			queue <- msg
		} else {
			go func() {
				time.Sleep(expire)
				queue <- msg
			}()
		}
	} else {
		return QueueNotExist
	}
	return
}

func (mq *LocalMQ) RegistryMessageHandler(topics []string, handler func(msg interface{})) (error error) {
	for _, topic := range topics {
		if queue, ok := mq.queues[topic]; ok && queue != nil {
			go func() {
				for msg := range queue {
					mq.handleError(topic, msg, handler)
				}
			}()
			return
		}
	}
	return QueueNotExist
}
