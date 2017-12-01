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
	l3Rules *packet.L3Rules

	outport uint
	inport  uint
)

// Main function for constructing packet processing graph.
func main() {
	flag.UintVar(&outport, "outport", 0, "port for sender")
	flag.UintVar(&inport, "inport", 0, "port for receiver")
	flag.Parse()

	// Init YANFF system at 16 available cores.
	config := flow.Config{
		CPUList: "0-15",
	}
	err := flow.SystemInit(&config)
	if err != nil {
		log.Fatal(err)
	}

	// Get splitting rules from access control file.
	l3Rules, err = packet.GetL3ACLFromORIG("test-handle2-l3rules.conf")
	if err != nil {
		log.Fatal(err)
	}

	// Receive packets from 0 port
	flow1, err := flow.SetReceiver(uint8(inport))
	if err != nil {
		log.Fatal(err)
	}

	// Handle packet flow
	err = flow.SetHandler(flow1, l3Handler, nil) // ~33% of packets should left in flow1
	if err != nil {
		log.Fatal(err)
	}

	// Send each flow to corresponding port. Send queues will be added automatically.
	err = flow.SetSender(flow1, uint8(outport))
	if err != nil {
		log.Fatal(err)
	}

	// Begin to process packets.
	err = flow.SystemStart()
	if err != nil {
		log.Fatal(err)
	}
}

func l3Handler(pkt *packet.Packet, context flow.UserContext) bool {
	return pkt.L3ACLPermit(l3Rules)
}
