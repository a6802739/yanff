package main

import "time"
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

func updateSeparateRules() {
	for true {
		time.Sleep(time.Second * 5)
		L3Rules = rules.GetL3RulesFromORIG("rules2.conf")
	}
}
