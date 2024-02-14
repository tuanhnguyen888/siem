package config

import (
	"context"
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/tuanhnguyen888/siem/collector/common"
	aws_s3sqs "github.com/tuanhnguyen888/siem/collector/input"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"sync"
)


func initCollection() ( *mongo.Collection, error) {
	mongoUrl, err := common.GoDotEnvVariable("SERVER_MONGODB")
	if err != nil {
		return nil , err
	}
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(mongoUrl))
	if err != nil {

		return nil, err
	}
	collection := client.Database("admin").Collection("log_source")
	return collection, nil
}

func InitInputConfig () ([]map[string]interface{},error) {
	collection , err := initCollection()
	if err != nil {
		return nil, err
	}

	filter := bson.D{}
	// Tìm tất cả các documents trong collection
	cursor, err := collection.Find(context.Background(), filter)
	if err != nil {
		return nil, fmt.Errorf("find data log source fail %v",err)
	}
	var data []map[string]interface{}
	for cursor.Next(context.Background()) {
		// Tạo một map để lưu trữ dữ liệu từ mỗi document
		var result map[string]interface{}
		if err = cursor.Decode(&result); err != nil {
			logrus.Error(err)
		}

		// Thêm map vào slice
		data = append(data, result["properties"].(map[string]interface{}))
	}
	if err = cursor.Err(); err != nil {
		return nil, err
	}

	return data,nil
}

func StartInputs(ctx context.Context, inputs []map[string]interface{},  msgChan chan<- []map[string]interface{}) error {
	var waitGroup sync.WaitGroup
	for _, input := range inputs {
		waitGroup.Add(1)
		fmt.Println(input)
		IInput, err := aws_s3sqs.InitHandler(ctx,input)
		if err != nil {
			logrus.Println("InitHandler fail", err)
			return err
		}
		go func() {
			defer func() {
				waitGroup.Done()
			}()

			err = IInput.Start(ctx,msgChan)
			if err != nil {
				logrus.Error(err)
			}

		}()
	}
	waitGroup.Wait()
	return nil
}