package config

import (
	"context"
	"correl/alertingSystem"
	"correl/input"
	"correl/threatAnalysis"
	"fmt"
	"github.com/sirupsen/logrus"
	"golang.org/x/sync/errgroup"
	"runtime"
	"sync"
)

type MsgChan chan map[string]interface{}

type Config struct {
	ChannelSizeInput  int
	ChannelSizeOutput int

	InputConfig   input.IInput
	FilterConfig  threatAnalysis.IThreatAnalysis
	OutputConfig  alertingSystem.IAlertSystem


	outputWorkers   int
	pipelineWorkers int

	chInFilter   MsgChan
	chFilterOut  MsgChan
	chOutDebug   MsgChan
	//RateCounters map[string]*ratecounter.RateCounter

	cancel    context.CancelFunc
	ctx       context.Context
	waitGroup *sync.WaitGroup
	eg        *errgroup.Group

}

func DefaultConfig() *Config {
	return &Config{
		ChannelSizeInput:  10000,
		ChannelSizeOutput: 10000,
		outputWorkers:     runtime.NumCPU(),
		pipelineWorkers:   runtime.NumCPU(),
	}
}

func (t *Config )Init() (err error) {
	t.initContext()
	t.chInFilter = make(MsgChan, t.ChannelSizeInput)
	t.chFilterOut = make(MsgChan, t.ChannelSizeOutput)


	t.InputConfig, err = input.InitInputConfig(t.ctx, t.chInFilter)
	if err != nil {
		return fmt.Errorf("init Input Config Fail: %v",err)
	}

	t.FilterConfig, err = threatAnalysis.InitThreatAnalysisConfig(t.ctx,t.chInFilter,t.chFilterOut)
	if err != nil {
		return fmt.Errorf("init Threat Analysis Config Fail: %v",err)
	}

	t.OutputConfig,err = alertingSystem.InitAlertSystem(t.ctx,t.chFilterOut)
	if err != nil {
		return fmt.Errorf("init Alert System Config Fail: %v",err)
	}

	return nil
}

func (t *Config) initContext() {
	t.waitGroup = &sync.WaitGroup{}

	ctx := context.Background()
	t.ctx, t.cancel = context.WithCancel(ctx)
	//ctx, _ = WithOSSignal(ctx, os.Interrupt, os.Kill)
	//t.eg, t.ctx = errgroup.WithContext(ctx)
}


func (t *Config) Stop() {
	if t.cancel != nil {
		logrus.Info("context cancel has been called")
		t.cancel()
	}

	t.Wait()
}

func (t *Config) Wait() {
	t.waitGroup.Wait()
}

func (t *Config) SendToInputFitlerChan(msg map[string]interface{}) {
	t.chInFilter <- msg
}

func (t *Config) SendToFilterOutputChan(msg map[string]interface{}) {
	t.chFilterOut <- msg
}

//// WithOSSignal tạo một context mới và hủy bỏ khi nhận được tín hiệu OS (Interrupt hoặc Kill).
//func WithOSSignal(parent context.Context, signals ...os.Signal) (context.Context, context.CancelFunc) {
//	ctx, cancel := context.WithCancel(parent)
//	sigCh := make(chan os.Signal, 1)
//	signal.Notify(sigCh, signals...)
//
//	var wg sync.WaitGroup
//	wg.Add(1)
//
//	go func() {
//		select {
//		case <-ctx.Done():
//		case sig := <-sigCh:
//			cancel()
//			// Xử lý tín hiệu nếu cần thiết
//			// Ví dụ: log.Printf("Received signal: %v", sig)
//		}
//		wg.Done()
//	}()
//
//	return ctx, func() {
//		cancel()
//		// Chờ goroutine xử lý tín hiệu kết thúc
//		wg.Wait()
//		signal.Stop(sigCh)
//		close(sigCh)
//	}
//}