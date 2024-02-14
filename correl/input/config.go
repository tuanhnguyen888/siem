package input

import (
	"context"
	"correl/common"
	"errors"
	"fmt"
	"github.com/IBM/sarama"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

)
type IInput interface {
	ListenForMessage( ) error

}

type InputConfig struct {
	BrokerList []string
	Topic string
	ConsumerGroup string
	ctx context.Context
	msg chan<- map[string]interface{}
	producer sarama.SyncProducer
	consumer sarama.Consumer
	collection  *mongo.Collection
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

func InitInputConfig(ctx context.Context,chanInput  chan<- map[string]interface{}) (IInput, error) {
	brokerList := getBrokerList()
	topic, err := common.GoDotEnvVariable("TOPIC_INPUT")
	if err != nil {
		return &InputConfig{} ,errors.New("no find .env key")
	}

	mongoUrl, err := common.GoDotEnvVariable("SERVER_MONGODB")
	if err != nil {
		return nil, err
	}
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(mongoUrl))
	if err != nil {
		return nil, errors.New("fail connect to mongo db :" + err.Error())
	}

	collection := client.Database("admin").Collection("event")

	return &InputConfig{
		BrokerList:    brokerList,
		Topic:         topic,
		ConsumerGroup: "siem",
		ctx: ctx,
		msg: chanInput,
		collection: collection,
		producer:      nil,
		consumer:      nil,
	}, nil
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