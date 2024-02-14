package config

import (
	"context"
	"github.com/sirupsen/logrus"
	aws_s3sqs "github.com/tuanhnguyen888/siem/collector/input"
	"testing"
	"time"
)

func TestInitInputConfig(t *testing.T) {
	data, _ := InitInputConfig()
	for _, input := range data {
		logrus.Println(input)
	}
}

func TestName(t *testing.T) {
	data, _ := InitInputConfig()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	for _, input := range data {
		testChan := make(chan []map[string]interface{}, 10)
		IInput, err := aws_s3sqs.InitHandler(ctx,input)
		if err != nil {
			logrus.Println("InitHandler fail", err)
		}
		go func() {
			err := IInput.Start(ctx,testChan)
			if err != nil {
				logrus.Error(err)
			}

		}()
	}
	time.Sleep(500000 * time.Millisecond)
	cancel()
	time.Sleep(500 * time.Millisecond)
}