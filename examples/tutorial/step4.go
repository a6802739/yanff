package main

import (
	"log"

	"github.com/intel-go/yanff/flow"
	"github.com/intel-go/yanff/packet"
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
	err = flow.SetSender(secondFlow, 1)
	if err != nil {
		log.Fatal(err)
	}

	err = flow.SystemStart()
	if err != nil {
		log.Fatal(err)
	}
}

func mySeparator(cur *packet.Packet, ctx flow.UserContext) bool {
	cur.ParseL3()
	if cur.GetIPv4() != nil {
		cur.ParseL4ForIPv4()
		if cur.GetTCPForIPv4() != nil && packet.SwapBytesUint16(cur.GetTCPForIPv4().DstPort) == 53 {
			return false
		}
	}
	return true
}
