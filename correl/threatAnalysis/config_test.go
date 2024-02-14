package threatAnalysis

import (
	"context"
	"fmt"
	"testing"
	"time"
)

func TestName(t *testing.T) {
	msgIn := make(chan map[string]interface{} , 1000)
	msgOut := make(chan map[string]interface{} , 1000)
	ctx, cancel := context.WithCancel(context.Background())
	conf, _ := InitThreatAnalysisConfig(ctx, msgIn,msgOut)
	defer cancel()

	go func() {
		for {
			time.Sleep(5 * time.Second)
			msg := map[string]interface{}{
				"url": "https://evil.com/",
			}
			msgIn <- msg
		}
	}()

	go func() {
		err := conf.ProcessEvent()
		if err != nil {
			fmt.Println(err)
		}
	}()

	time.Sleep(100 * time.Second)
	cancel()
	time.Sleep(200 * time.Millisecond)
}
