package config

import (
	"context"
	output2 "github.com/tuanhnguyen888/siem/collector/output"
	"sync"
)

type Config struct {
	Input []map[string]interface{}
	Output output2.IOutput
	ChFilterOut chan []map[string]interface{}
	channelSizeInput  int

	Cancel    context.CancelFunc
	Ctx       context.Context
	WaitGroup *sync.WaitGroup
}

func DefaultConfig() *Config {
	return &Config{
		channelSizeInput: 10000,
	}
}

func(t *Config) Init() error{
	t.initContext()
	t.ChFilterOut = make(chan []map[string]interface{},t.channelSizeInput)

	inputData, err := InitInputConfig()
	if err != nil {
		return err
	}
	t.Input = inputData

	output, err := output2.InitInputConfig(t.Ctx,t.ChFilterOut)
	if err != nil {
		return err
	}
	t.Output = output


	return nil
	
}

func (t *Config) initContext() {
	t.WaitGroup = &sync.WaitGroup{}

	ctx := context.Background()
	t.Ctx, t.Cancel = context.WithCancel(ctx)
}

