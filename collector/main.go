package main

import (
	"github.com/tuanhnguyen888/siem/collector/cmd"
	"log"
)

func main()  {
	if err := cmd.Run(); err != nil {
		log.Panicln(err)
	}
}