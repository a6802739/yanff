package main

import (
	"log"
	"time"

	"github.com/intel-go/yanff/flow"
	"github.com/intel-go/yanff/packet"
)

var (
	l3Rules *packet.L3Rules
)

func main() {
	config := flow.Config{}
	err := flow.SystemInit(&config)
	if err != nil {
		log.Fatal(err)
	}

	initCommonState()

	l3Rules, err = packet.GetL3ACLFromORIG("rules1.conf")
	if err != nil {
		log.Fatal(err)
	}

	go updateSeparateRules()
	firstFlow, err := flow.SetReceiver(0)
	if err != nil {
		log.Fatal(err)
	}
	secondFlow, err := flow.SetSeparator(firstFlow, mySeparator, nil)
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
	err = flow.SetStopper(secondFlow)
	if err != nil {
		log.Fatal(err)
	}
	err = flow.SystemStart()
	if err != nil {
		log.Fatal(err)
	}
}

func mySeparator(cur *packet.Packet, ctx flow.UserContext) bool {
	localL3Rules := l3Rules
	return cur.L3ACLPermit(localL3Rules)
}

func updateSeparateRules() {
	for {
		time.Sleep(time.Second * 5)
		var err error
		l3Rules, err = packet.GetL3ACLFromORIG("rules1.conf")
		if err != nil {
			log.Fatal(err)
		}
	}
}
