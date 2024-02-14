package output

import (
	"encoding/json"
	"github.com/IBM/sarama"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

func (o *Output)PushEvtToKafka() error  {
	logrus.Println("Collector pushing event to kafka log raw ..........")
	for  {
		select {
		case logs, ok := <-o.chFilter:
			if !ok {
				return errors.New("channel input filter closed. Exiting")
			}
			for _ , log := range logs{
				jsonData, err := json.Marshal(log)
				if err != nil {
					logrus.Println("Error Marshal to []byte:", err)
				}

				msg := &sarama.ProducerMessage{
					Topic:     o.Topic,
					Value:     sarama.ByteEncoder(jsonData),
				}

				_, _, err = o.producer.SendMessage(msg)
				if err != nil {
					logrus.Errorf("Failed to produce message: %v", err)
				}
			}

		}
	}
}