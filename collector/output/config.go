package output

import (
	"context"
	"fmt"
	"github.com/IBM/sarama"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/tuanhnguyen888/siem/collector/common"
)

type IOutput interface {
	PushEvtToKafka() error
}

type Output struct {
	BrokerList []string
	Topic string
	ConsumerGroup string
	ctx context.Context
	chFilter <-chan []map[string]interface{}
	producer sarama.SyncProducer
}

func InitInputConfig(ctx context.Context, msg <-chan []map[string]interface{} ) ( IOutput ,error ) {
	brokerList := getBrokerList()
	topic, err := common.GoDotEnvVariable("TOPIC_INPUT")
	if err != nil {
		return &Output{} ,errors.New("no find .env key")
	}
	produce, err := initPubKafka()
	if err != nil {
		logrus.Errorf("Init produce fail: %v", err)
		return &Output{} ,err
	}


	return &Output{
		BrokerList:    brokerList,
		Topic:         topic,
		ConsumerGroup: "siem",
		ctx: ctx,
		chFilter: msg,
		producer:      produce,
	}, nil

}

func initPubKafka() (sarama.SyncProducer,error) {
	// Get the number of broker
	brks := getBrokerList()

	config := sarama.NewConfig()
	config.Producer.RequiredAcks = sarama.WaitForLocal
	config.Producer.Return.Successes = true
	producer, err := sarama.NewSyncProducer(brks, config)
	if err != nil {
		return nil, err
	}

	return producer, nil
}

func getBrokerList() (brks []string)  {
	numDbs := 1
	for {
		brk, err := common.GoDotEnvVariable(fmt.Sprintf("SERVER_BROKER%v", numDbs))
		if err != nil {
			break
		}
		brks = append(brks,brk)
		numDbs++
	}
	return brks
}

