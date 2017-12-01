// Copyright 2017 Intel Corporation.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"flag"
	"log"
	"github.com/intel-go/yanff/flow"
)

var (
	outport uint
	inport  uint
	cores   string
)

// Main function for constructing packet processing graph.
func main() {
	flag.UintVar(&outport, "outport", 1, "port for sender")
	flag.UintVar(&inport, "inport", 0, "port for receiver")
	flag.StringVar(&cores, "cores", "0-15", "Specifies CPU cores to be used by YANFF library")

	// Initialize YANFF library at requested cores.
	config := flow.Config{
		CPUList: cores,
	}
	err := flow.SystemInit(&config)
	if err != nil {
		log.Fatal(err)
	}

	// Receive packets from 0 port and send to 1 port.
	flow1, err := flow.SetReceiver(uint8(inport))
	if err != nil {
		log.Fatal(err)
	}
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
