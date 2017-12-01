package main

import (
	"log"

	"github.com/intel-go/yanff/flow"
)

func main() {
	// Init YANFF system
	config := flow.Config{}
	err := flow.SystemInit(&config)
	if err != nil {
		log.Fatal(err)
	}
	initCommonState()

	err = flow.SystemStart()
	if err != nil {
		log.Fatal(err)
	}
}
