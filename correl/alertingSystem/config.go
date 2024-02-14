package alertingSystem

import (
	"context"
	"correl/common"
	"errors"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type IAlertSystem interface {
	PushEventToMongo() error
}

type AlertSystem struct {
	msgInput <-chan map[string]interface{}
	ctx context.Context
	client *mongo.Client
	collection  *mongo.Collection
}

func InitAlertSystem(ctx context.Context, msgChan <-chan map[string]interface{}) (IAlertSystem,error) {
	mongoUrl, err := common.GoDotEnvVariable("SERVER_MONGODB")
	if err != nil {
		return nil, err
	}
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(mongoUrl))
	if err != nil {
		return nil, errors.New("fail connect to mongo db :" + err.Error())
	}
	collection := client.Database("admin").Collection("alert")

	return &AlertSystem{
		msgInput: msgChan,
		ctx:      ctx,
		client:   client,
		collection: collection,
	}, nil
}

func (f *AlertSystem) PushEventToMongo () error {
	logrus.Println("Pushing event to mongo ..........")
	for  {
		select {
		case alert, ok := <-f.msgInput:
			if !ok {
				return errors.New("channel input filter closed. Exiting")
			}
			_ , err := f.collection.InsertOne(context.Background(),alert)
			if err != nil{
				logrus.Errorf("Insert alert %s fail", alert["id"].(string))
			}

			logrus.Printf("Push success alert_id: %s to mongoDB", alert["id"].(string))
		}
	}
}

