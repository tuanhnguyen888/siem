package cmd

import (
	"github.com/sirupsen/logrus"
	"testing"
)

func TestRun(t *testing.T) {
	err := Run()
	if err != nil {
		logrus.Println(err)
	}


}