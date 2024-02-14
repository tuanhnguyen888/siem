package input

import (
	"context"
	"correl/common"
	"encoding/json"
	"fmt"
	"github.com/IBM/sarama"
	"github.com/hashicorp/go-uuid"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"testing"

	"time"
)

func TestInitKafka(t *testing.T) {
	produce, err := initPubKafka()
	fmt.Println(produce,"--",err)
}

func TestPutEventKafka(t *testing.T){
	produce, err := initPubKafka()
	if err != nil {
		logrus.Errorf("Init produce fail: %v", err)
		return
	}
	id, _ := uuid.GenerateUUID()
	data := map[string]interface{}{
		"id": id,
		"description_ip": "",
		"source_ip":"",
	}

	// Chuyển đổi map thành slice []byte
	jsonData, err := json.Marshal(data)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}


	msg := &sarama.ProducerMessage{
		Topic:     "log_raw",
		Value:     sarama.ByteEncoder(jsonData),
	}

	partition, offset, err := produce.SendMessage(msg)
	if err != nil {
		log.Fatalf("Failed to produce message: %v", err)
	}

	fmt.Printf("Produced message to partition %d with offset %d\n", partition, offset)

}

func TestPushMongo(t *testing.T) {
	//id, _ := uuid.GenerateUUID()
	//data := map[string]interface{}{
	//	"id": id,
	//	"input_name": "Amazon Web Services",
	//	"protocol": "s3_sqs",
	//	"properties": map[string]interface{}{
	//		"queue_url" : "https://sqs.us-east-1.amazonaws.com/887134122148/vpc-anhnt",
	//		"log_type" : "vpc_flow_logs",
	//		"visibility_timeout"  : int64(300),
	//		"wait_time" : int64(120),
	//		"max_number_of_messages": int64(10),
	//		"api_timeout" : int64(120),
	//		"access_key_id": "AKIA45DKQOCSDRHEU5Q7",
	//		"secret_access_key": "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",
	//		"proxy_url" : "",
	//	},
	//	"created_at": int64(1703956316000),
	//	"updated_at" : int64(1703956881000),
	//}
	//bytes, _ := bcrypt.GenerateFromPassword([]byte("123456789"), 14)
	dataUser := map[string]interface{}{
		"input_name": "Amazon Web Services",
		"protocol": "s3_sqs",
		"data": map[string]interface{}{
			"queue url" : map[string]interface{}{
				"key": "queue_url",
				"type": "string",
				"required": true,
			},
			"log type": map[string]interface{}{
				"key": "log_type",
				"type": "option string",
				"format": []string{"cloudtrail", "vpc_flow_log","route53"},
				"required": true,
			},
			"visibility timeout": map[string]interface{}{
					"key": "visibility_timeout",
					"type": "number",
					"required": true,
				},
		},
		"created_at": int64(1703956316000),
		"updated_at" : int64(1703956881000),
	}
	mongoUrl, err := common.GoDotEnvVariable("SERVER_MONGODB")
	if err != nil {
		logrus.Errorf("go dot %v",err)
		return
	}
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(mongoUrl))
	if err != nil {
		logrus.Error(err)
		return
	}
	collection := client.Database("admin").Collection("prototype")
	_ , err = collection.InsertOne(context.Background(),dataUser)
	if err != nil {
		logrus.Error("InsertOne fail", err)
		return
	}

	logrus.Println("success")
}

func TestListenEvent(t *testing.T) {

	msg := make(chan map[string]interface{} , 1000)
	ctx, cancel := context.WithCancel(context.Background())
	input, _ := InitInputConfig(ctx, msg)
	defer cancel()
	go func() {
		err :=  input.ListenForMessage()
		fmt.Println(err)
	}()

	time.Sleep(1221 * time.Second)
	cancel()
	time.Sleep(200 * time.Millisecond)

}