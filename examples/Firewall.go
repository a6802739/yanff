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
	flag.UintVar(&outport, "outport", 1, "port for sender")
	flag.UintVar(&inport, "inport", 0, "port for receiver")
	flag.Parse()

	// Initialize YANFF library at 8 cores by default
	config := flow.Config{
		CPUList: "0-7",
	}
	err := flow.SystemInit(&config)
	if err != nil {
		log.Fatal(err)
	}

	// Get filtering rules from access control file.
	l3Rules, err = packet.GetL3ACLFromORIG("Firewall.conf")
	if err != nil {
		log.Fatal(err)
	}

	// Receive packets from zero port. Receive queue will be added automatically.
	inputFlow, err := flow.SetReceiver(uint8(inport))
	if err != nil {
		log.Fatal(err)
	}

	// Separate packet flow based on ACL.
	rejectFlow, err := flow.SetSeparator(inputFlow, l3Separator, nil)
	if err != nil {
		log.Fatal(err)
	}

	// Drop rejected packets.
	err = flow.SetStopper(rejectFlow)
	if err != nil {
		log.Fatal(err)
	}

	// Send accepted packets to first port. Send queue will be added automatically.
	err = flow.SetSender(inputFlow, uint8(outport))
	if err != nil {
		log.Fatal(err)
	}

	// Begin to process packets.
	err = flow.SystemStart()
	if err != nil {
		log.Fatal(err)
	}
}

// User defined function for separating packets
func l3Separator(currentPacket *packet.Packet, context flow.UserContext) bool {
	// Return whether packet is accepted or not. Based on ACL rules.
	return currentPacket.L3ACLPermit(l3Rules)
}
