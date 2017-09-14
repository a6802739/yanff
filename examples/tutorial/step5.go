package main

import "github.com/intel-go/yanff/flow"
import "github.com/intel-go/yanff/packet"
import "github.com/intel-go/yanff/rules"

var (
	L3Rules *rules.L3Rules
)

func main() {
	config := flow.Config{}
	flow.SystemInit(&config)
	initCommonState()
	L3Rules = rules.GetL3RulesFromJSON("step5.json")
	firstFlow := flow.SetReceiver(0)
	secondFlow := flow.SetSeparator(firstFlow, mySeparator, nil)
	flow.SetHandler(firstFlow, modifyPacket[0], nil)
	flow.SetHandler(secondFlow, modifyPacket[1], nil)
	flow.SetSender(firstFlow, 0)
	flow.SetSender(secondFlow, 1)

	flow.SystemStart()
}

func mySeparator(cur *packet.Packet, ctx flow.UserContext) bool {
	cur.ParseIPv4TCP()
	return rules.L3ACLPermit(cur, L3Rules)
}
