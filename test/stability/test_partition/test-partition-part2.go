// Copyright 2017 Intel Corporation.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"flag"
	"log"

	"github.com/intel-go/yanff/flow"
	"github.com/intel-go/yanff/packet"
	"github.com/intel-go/yanff/test/stability/stabilityCommon"
)

var (
	inport   uint
	outport1 uint
	outport2 uint

	fixMACAddrs1 func(*packet.Packet, flow.UserContext)
	fixMACAddrs2 func(*packet.Packet, flow.UserContext)
)

// Main function for constructing packet processing graph.
func main() {
	flag.UintVar(&inport, "inport", 0, "port for receiver")
	flag.UintVar(&outport1, "outport1", 0, "port for 1st sender")
	flag.UintVar(&outport2, "outport2", 1, "port for 2nd sender")
	configFile := flag.String("config", "config.json", "Specify config file name")
	target := flag.String("target", "nntsat01g4", "Target host name from config file")
	flag.Parse()

	// Init YANFF system at 16 available cores.
	config := flow.Config{
		CPUList: "0-15",
	}
	err := flow.SystemInit(&config)
	if err != nil {
		log.Fatal(err)
	}
	stabilityCommon.InitCommonState(*configFile, *target)
	fixMACAddrs1 = stabilityCommon.ModifyPacket[outport1].(func(*packet.Packet, flow.UserContext))
	fixMACAddrs2 = stabilityCommon.ModifyPacket[outport2].(func(*packet.Packet, flow.UserContext))

	// Receive packets from 0 port
	flow1, err := flow.SetReceiver(uint8(inport))
	if err != nil {
		log.Fatal(err)
	}
	flow2, err := flow.SetPartitioner(flow1, 1000, 100)
	if err != nil {
		log.Fatal(err)
	}

	err = flow.SetHandler(flow1, fixPackets1, nil)
	if err != nil {
		log.Fatal(err)
	}
	err = flow.SetHandler(flow2, fixPackets2, nil)
	if err != nil {
		log.Fatal(err)
	}

	err = flow.SetSender(flow1, uint8(outport1))
	if err != nil {
		log.Fatal(err)
	}
	err = flow.SetSender(flow2, uint8(outport2))
	if err != nil {
		log.Fatal(err)
	}

	// Begin to process packets.
	err = flow.SystemStart()
	if err != nil {
		log.Fatal(err)
	}
}

func fixPackets1(pkt *packet.Packet, ctx flow.UserContext) {
	if stabilityCommon.ShouldBeSkipped(pkt) {
		return
	}
	fixMACAddrs1(pkt, ctx)
}

func fixPackets2(pkt *packet.Packet, ctx flow.UserContext) {
	if stabilityCommon.ShouldBeSkipped(pkt) {
		return
	}
	fixMACAddrs2(pkt, ctx)
}
