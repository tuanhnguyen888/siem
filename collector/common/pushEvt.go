package common

import "github.com/sirupsen/logrus"

func PushEvtToChan(evts []map[string]interface{} , msgChan chan<- []map[string]interface{}) {
	msgChan <- evts
	logrus.Printf("Done push event to chan")
}

