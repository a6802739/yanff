package main

import (
	"log"

	"github.com/intel-go/yanff/flow"
)

func main() {
	config := flow.Config{}
	err := flow.SystemInit(&config)
	if err != nil {
		log.Fatal(err)
	}

	initCommonState()

	firstFlow, err := flow.SetReceiver(0)
	if err != nil {
		log.Fatal(err)
	}
	secondFlow, err := flow.SetPartitioner(firstFlow, 300, 300)
	if err != nil {
		log.Fatal(err)
	}
	err = flow.SetHandler(firstFlow, modifyPacket[0], nil)
	if err != nil {
		log.Fatal(err)
	}
	err = flow.SetHandler(secondFlow, modifyPacket[1], nil)
	if err != nil {
		log.Fatal(err)
	}
	err = flow.SetSender(firstFlow, 0)
	if err != nil {
		log.Fatal(err)
	}
	err = flow.SetSender(secondFlow, 1)

	err = flow.SystemStart()
	if err != nil {
		log.Fatal(err)
	}
}
