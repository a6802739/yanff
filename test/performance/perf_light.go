// Copyright 2017 Intel Corporation.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"flag"
	"log"

	"github.com/intel-go/yanff/flow"
	"github.com/intel-go/yanff/packet"
)

var (
	mode uint

	inport1     uint
	inport2     uint
	outport1    uint
	outport2    uint
	noscheduler bool
)

func main() {
	flag.UintVar(&mode, "mode", 0, "Benching mode: 0 - empty, 1 - parsing, 2 - parsing, reading, writing")
	flag.UintVar(&outport1, "outport1", 1, "port for 1st sender")
	flag.UintVar(&outport2, "outport2", 1, "port for 2nd sender")
	flag.UintVar(&inport1, "inport1", 0, "port for 1st receiver")
	flag.UintVar(&inport2, "inport2", 0, "port for 2nd receiver")
	flag.BoolVar(&noscheduler, "no-scheduler", false, "disable scheduler")
	flag.Parse()

	// Initialize YANFF library at 15 cores by default
	config := flow.Config{
		CPUList:          "0-14",
		DisableScheduler: noscheduler,
	}
	err := flow.SystemInit(&config)
	if err != nil {
		log.Fatal(err)
	}

	// Receive packets from zero port. One queue per receive will be added automatically.
	firstFlow0, err := flow.SetReceiver(uint8(inport1))
	if err != nil {
		log.Fatal(err)
	}
	firstFlow1, err := flow.SetReceiver(uint8(inport2))
	if err != nil {
		log.Fatal(err)
	}

	firstFlow, err := flow.SetMerger(firstFlow0, firstFlow1)
	if err != nil {
		log.Fatal(err)
	}

	// Handle second flow via some heavy function
	if mode == 0 {
		err = flow.SetHandler(firstFlow, heavyFunc0, nil)
		if err != nil {
			log.Fatal(err)
		}
	} else if mode == 1 {
		err = flow.SetHandler(firstFlow, heavyFunc1, nil)
		if err != nil {
			log.Fatal(err)
		}
	} else {
		err = flow.SetHandler(firstFlow, heavyFunc2, nil)
		if err != nil {
			log.Fatal(err)
		}
	}

	// Split for two senders and send
	secondFlow, err := flow.SetPartitioner(firstFlow, 150, 150)
	if err != nil {
		log.Fatal(err)
	}

	err = flow.SetSender(firstFlow, uint8(outport1))
	if err != nil {
		log.Fatal(err)
	}
	err = flow.SetSender(secondFlow, uint8(outport2))
	if err != nil {
		log.Fatal(err)
	}

	err = flow.SystemStart()
	if err != nil {
		log.Fatal(err)
	}
}

func heavyFunc0(currentPacket *packet.Packet, context flow.UserContext) {
}

func heavyFunc1(currentPacket *packet.Packet, context flow.UserContext) {
	currentPacket.ParseL3()
}

func heavyFunc2(currentPacket *packet.Packet, context flow.UserContext) {
	currentPacket.ParseL3()
	ipv4 := currentPacket.GetIPv4()
	if ipv4 != nil {
		T := ipv4.DstAddr
		ipv4.SrcAddr = 263 + (T)
	}
}
