package cmd

import (
	"correl/config"
	"fmt"
	"github.com/sirupsen/logrus"
	"sync"
)

func Run()  {
	conf := config.DefaultConfig()
	conf.Init()


	var wg sync.WaitGroup
	numGoroutines := 3
	wg.Add(numGoroutines)
	for i := 0; i < numGoroutines; i++ {
		go func(index int) {
			// Gọi hàm tương ứng với goroutine
			switch index {
			case 0:
				err := conf.InputConfig.ListenForMessage()
				if err != nil {
					logrus.Errorf(" ListenForMessage fail: %v", err)
				}

			case 1:
				err := conf.FilterConfig.AnalysisEvent()
				if err != nil {
					logrus.Errorf(" AnalysisEvent fail: %v", err)
				}
			case 2:
				err := conf.OutputConfig.PushEventToMongo()
				if err != nil {
					logrus.Errorf(" PushEventToMongo fail: %v", err)
				}
			}

			// Đánh dấu goroutine đã hoàn thành công việc
			wg.Done()
		}(i)
	}
	// Đợi tất cả các goroutine hoàn thành
	wg.Wait()

	// Tất cả các goroutine đã hoàn thành, chương trình sẽ kết thúc
	fmt.Println("All goroutines have closed")

}