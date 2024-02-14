package input

import (
	"context"
	"correl/common"
	"encoding/json"
	"fmt"
	"github.com/IBM/sarama"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

func (i *InputConfig ) ListenForMessage() error {
	brks := getBrokerList()
	logrus.Info("Kafka server",brks)
	config := sarama.NewConfig()
	config.Consumer.Return.Errors = true

	// Tạo đối tượng consumer
	consumer, err := sarama.NewConsumerGroup(brks, i.ConsumerGroup , config)
	if err != nil {
		fmt.Printf("Error creating consumer group: %v\n", err)
		return err
	}

	logrus.Info("listening for event from kafka")
	go func() {
		for {
			select {
			case <-i.ctx.Done():
				break
			default:
				handler := &ConsumerHandler{
					Msg : i.msg,
					collection: i.collection,
				}
				if err = consumer.Consume(i.ctx, []string{i.Topic}, handler); err != nil {
					logrus.Error( "Consumer error: " ,err)
					time.Sleep(time.Second)
				}
			}
		}
	}()
	<- i.ctx.Done()
	return nil
}

func createEvt(msg *sarama.ConsumerMessage)  (map[string]interface{}, error) {
	var evt map[string]interface{}

	err := json.Unmarshal( msg.Value , &evt)
	if err != nil {
		logrus.Error("Error Unmarshal:", err)
		return evt, err
	}

	return evt, nil
}

//consumer
type ConsumerHandler struct{
	Msg chan<- map[string]interface{}
	collection  *mongo.Collection
}

func (h *ConsumerHandler) Setup(sarama.ConsumerGroupSession) error {
	// Mark the consumer as ready
	//consumer.ready <- true
	return nil
}

// Cleanup thực hiện các công việc dọn dẹp sau khi consumer dừng hoạt động
func (h *ConsumerHandler) Cleanup(_ sarama.ConsumerGroupSession) error {
	return nil
}

// ConsumeClaim xử lý các thông điệp từ Kafka
func (h *ConsumerHandler) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {

	for msg := range claim.Messages() {
		evt, err := createEvt(msg)
		if err != nil {
			logrus.Error(err)
		}
		dataMongo := createEvtMongo(evt)
		_ , err = h.collection.InsertOne(context.Background(),dataMongo)
		if err != nil{
			logrus.Errorf("Insert event fail, %v", err)
		}

		common.PushEvtToChan(evt, h.Msg)
		logrus.Printf("Received event from kafka: %v ", evt["id"].(string))

		// Đánh dấu thông điệp đã được xử lý
		session.MarkMessage(msg, "")
	}
	return nil
}

func createEvtMongo(evt map[string]interface{})  map[string]interface{} {
	return map[string]interface{}{
		"log_source": evt["provider"].(string),
		"data" : evt,
		"local_timestamp" : evt["local_timestamp"],
		"timestamp": time.Now().UnixMilli(),
	}
}