package main

import "time"
import "github.com/intel-go/yanff/common"
import "github.com/intel-go/yanff/flow"
import "github.com/intel-go/yanff/packet"
import "github.com/intel-go/yanff/rules"

var (
	L3Rules *rules.L3Rules
)
const flowN = 3

func main() {
	config := flow.Config{}
	flow.SystemInit(&config)
	initCommonState()
	L3Rules = rules.GetL3RulesFromORIG("rules2.conf")
	go updateSeparateRules()
	firstFlow := flow.SetReceiver(0)
	outputFlows := flow.SetSplitter(firstFlow, mySplitter, flowN, nil)
	flow.SetStopper(outputFlows[0])
	flow.SetHandler(outputFlows[1], myHandler, nil)
	for i := 1; i < flowN; i++ {
		flow.SetHandler(outputFlows[i], modifyPacket[i-1], nil)
		flow.SetSender(outputFlows[i], uint8(i-1))
	}
	flow.SystemStart()
}

func mySplitter(cur *packet.Packet, ctx flow.UserContext) uint {
	cur.ParseIPv4TCP()
	localL3Rules := L3Rules
	return rules.L3ACLPort(cur, localL3Rules)
}

func myHandler(cur *packet.Packet, ctx flow.UserContext) {
	cur.EncapsulateHead(common.EtherLen, common.IPv4MinLen)
	cur.ParseIPv4()
	cur.IPv4.SrcAddr = packet.IPv4(111, 22, 3, 0)
	cur.IPv4.DstAddr = packet.IPv4(3, 22, 111, 0)
	cur.IPv4.VersionIhl = 0x45
	cur.IPv4.NextProtoID = 0x04
}

func updateSeparateRules() {
	for true {
		time.Sleep(time.Second * 5)
		L3Rules = rules.GetL3RulesFromORIG("rules2.conf")
	}
}
