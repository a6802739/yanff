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
	inport1 uint
	inport2 uint
	outport uint

	fixMACAddrs func(*packet.Packet, flow.UserContext)
)

// Main function for constructing packet processing graph.
func main() {
	flag.UintVar(&inport1, "inport1", 0, "port for 1st receiver")
	flag.UintVar(&inport2, "inport2", 1, "port for 2nd receiver")
	flag.UintVar(&outport, "outport", 0, "port for sender")
	configFile := flag.String("config", "config.json", "Specify config file name")
	target := flag.String("target", "nntsat01g4", "Target host name from config file")
	flag.Parse()

	// Init YANFF system at requested number of cores.
	config := flow.Config{
		CPUList: "0-15",
	}
	err := flow.SystemInit(&config)
	if err != nil {
		log.Fatal(err)
	}
	stabilityCommon.InitCommonState(*configFile, *target)
	fixMACAddrs = stabilityCommon.ModifyPacket[outport].(func(*packet.Packet, flow.UserContext))

	// Receive packets from 0 and 1 ports
	inputFlow1, err := flow.SetReceiver(uint8(inport1))
	if err != nil {
		log.Fatal(err)
	}
	inputFlow2, err := flow.SetReceiver(uint8(inport2))
	if err != nil {
		log.Fatal(err)
	}

	outputFlow, err := flow.SetMerger(inputFlow1, inputFlow2)
	if err != nil {
		log.Fatal(err)
	}
	err = flow.SetHandler(outputFlow, fixPackets, nil)
	if err != nil {
		log.Fatal(err)
	}
	err = flow.SetSender(outputFlow, uint8(outport))
	if err != nil {
		log.Fatal(err)
	}

	// Begin to process packets.
	err = flow.SystemStart()
	if err != nil {
		log.Fatal(err)
	}
}

func fixPackets(pkt *packet.Packet, ctx flow.UserContext) {
	if stabilityCommon.ShouldBeSkipped(pkt) {
		return
	}
	fixMACAddrs(pkt, ctx)
}
