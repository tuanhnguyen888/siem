package common

import (
	"github.com/sirupsen/logrus"
)

func PushEvtToChan(evt map[string]interface{} , msgChan chan<- map[string]interface{}) {
	val, ok := evt["id"].(string)
	if ok {
		msgChan <- evt
		logrus.Printf("Done push event %s to chan",val)
	}
	 return
}

