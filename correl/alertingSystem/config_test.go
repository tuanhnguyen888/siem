package alertingSystem

import (
	"context"
	"github.com/sirupsen/logrus"
	"testing"
)

func TestInitAlertSystem(t *testing.T) {
	msgChanlTest := make(chan map[string]interface{}, 1000)
	_, err := InitAlertSystem(context.Background(), msgChanlTest)
	logrus.Warn(err)
}