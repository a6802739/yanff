package main

import (
	"log"
	"time"

	"github.com/intel-go/yanff/common"
	"github.com/intel-go/yanff/flow"
	"github.com/intel-go/yanff/packet"
)

var (
	l3Rules *packet.L3Rules
)

const flowN = 3

func main() {
	config := flow.Config{}
	err := flow.SystemInit(&config)
	if err != nil {
		log.Fatal(err)
	}

	initCommonState()

	l3Rules, err = packet.GetL3ACLFromORIG("packet..conf")
	if err != nil {
		log.Fatal(err)
	}
	go updateSeparateRules()
	firstFlow, err := flow.SetReceiver(0)
	if err != nil {
		log.Fatal(err)
	}
	outputFlows, err := flow.SetSplitter(firstFlow, mySplitter, flowN, nil)
	if err != nil {
		log.Fatal(err)
	}

	err = flow.SetStopper(outputFlows[0])
	if err != nil {
		log.Fatal(err)
	}
	err = flow.SetHandler(outputFlows[1], myHandler, nil)
	if err != nil {
		log.Fatal(err)
	}
	for i := 1; i < flowN; i++ {
		err = flow.SetHandler(outputFlows[i], modifyPacket[i-1], nil)
		if err != nil {
			log.Fatal(err)
		}
		err = flow.SetSender(outputFlows[i], uint8(i-1))
		if err != nil {
			log.Fatal(err)
		}
	}
	err = flow.SystemStart()
	if err != nil {
		log.Fatal(err)
	}
}

func mySplitter(cur *packet.Packet, ctx flow.UserContext) uint {
	localL3Rules := l3Rules
	return cur.L3ACLPort(localL3Rules)
}

func myHandler(curV []*packet.Packet, num uint, ctx flow.UserContext) {
	for i := uint(0); i < num; i++ {
		cur := curV[i]
		cur.EncapsulateHead(common.EtherLen, common.IPv4MinLen)
		cur.ParseL3()
		cur.GetIPv4().SrcAddr = packet.BytesToIPv4(111, 22, 3, 0)
		cur.GetIPv4().DstAddr = packet.BytesToIPv4(3, 22, 111, 0)
		cur.GetIPv4().VersionIhl = 0x45
		cur.GetIPv4().NextProtoID = 0x04
	}
}

func updateSeparateRules() {
	for {
		time.Sleep(time.Second * 5)
		var err error
		l3Rules, err = packet.GetL3ACLFromORIG("rules2.conf")
		if err != nil {
			log.Fatal(err)
		}
	}
}
