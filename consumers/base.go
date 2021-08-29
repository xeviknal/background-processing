package consumers

import (
	"log"

	"github.com/xeviknal/background-processing/buffer"
)

type Consumers []Consumer
type Consumer interface {
	Perform(params map[string]interface{})
}

var consumerMap = map[string]Consumer{
	"JobsConsumer":         JobsConsumer{},
	"JobCancellerConsumer": JobCancellerConsumer{},
}

func ConsumeTasks() {
	queue := buffer.GetBuffer()
	for {
		task := *queue.Next()
		consumer, found := consumerMap[task.GetConsumer()]
		if !found {
			log.Fatalf("consumer not found: %s", task.GetConsumer())
			continue
		}
		go consumer.Perform(task.GetDescriptor())
	}
}
