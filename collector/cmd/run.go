package cmd

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/tuanhnguyen888/siem/collector/config"
)

func Run() error {
	conf := config.DefaultConfig()
	err := conf.Init()
	if err != nil {
		return err
	}
	conf.WaitGroup.Add(2)
	for i := 0; i < 2; i++ {
		go func(index int) {
			// Gọi hàm tương ứng với goroutine
			switch index {
			case 0:
				err := 	config.StartInputs(conf.Ctx,conf.Input,conf.ChFilterOut)
				if err != nil {
					logrus.Errorf(" StartInputs fail: %v", err)
				}
			case 1:
				err := conf.Output.PushEvtToKafka()
				if err != nil {
					logrus.Errorf(" PushEvtToKafka fail: %v", err)
				}
			}

			conf.WaitGroup.Done()
		}(i)
	}

	conf.WaitGroup.Wait()
	fmt.Println("All goroutines have closed")

	return nil
}
