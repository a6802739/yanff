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
	load uint
	mode uint

	inport1     uint
	inport2     uint
	outport1    uint
	outport2    uint
	noscheduler bool
)

func main() {
	flag.UintVar(&load, "load", 1000, "Use this for regulating 'load intensity', number of iterations")
	flag.UintVar(&mode, "mode", 2, "Benching mode: 2, 12 - two handles; 3, 13 - tree handles; 4, 14 - four handles. 2,3,4 - one flow; 12,13,14 - two flows")
	flag.UintVar(&outport1, "outport1", 1, "port for 1st sender")
	flag.UintVar(&outport2, "outport2", 1, "port for 2nd sender")
	flag.UintVar(&inport1, "inport1", 0, "port for 1st receiver")
	flag.UintVar(&inport2, "inport2", 0, "port for 2nd receiver")
	flag.BoolVar(&noscheduler, "no-scheduler", false, "disable scheduler")
	flag.Parse()

	// Initialize YANFF library at 35 cores by default
	config := flow.Config{
		CPUList:          "0-34",
		DisableScheduler: noscheduler,
	}
	err := flow.SystemInit(&config)
	if err != nil {
		log.Fatal(err)
	}

	var tempFlow *flow.Flow
	var afterFlow *flow.Flow

	// Receive packets from zero port. One queue will be added automatically.
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

	if mode > 10 {
		tempFlow, err = flow.SetPartitioner(firstFlow, 150, 150)
		if err != nil {
			log.Fatal(err)
		}
	}

	// Handle second flow via some heavy function
	err = flow.SetHandler(firstFlow, heavyFunc, nil)
	if err != nil {
		log.Fatal(err)
	}
	err = flow.SetHandler(firstFlow, heavyFunc, nil)
	if err != nil {
		log.Fatal(err)
	}
	if mode%10 > 2 {
		err = flow.SetHandler(firstFlow, heavyFunc, nil)
		if err != nil {
			log.Fatal(err)
		}
	}
	if mode%10 > 3 {
		err = flow.SetHandler(firstFlow, heavyFunc, nil)
		if err != nil {
			log.Fatal(err)
		}
	}
	if mode > 10 {
		err = flow.SetHandler(tempFlow, heavyFunc, nil)
		if err != nil {
			log.Fatal(err)
		}
		err = flow.SetHandler(tempFlow, heavyFunc, nil)
		if err != nil {
			log.Fatal(err)
		}
		if mode%10 > 2 {
			err = flow.SetHandler(tempFlow, heavyFunc, nil)
			if err != nil {
				log.Fatal(err)
			}
		}
		if mode%10 > 3 {
			err = flow.SetHandler(tempFlow, heavyFunc, nil)
			if err != nil {
				log.Fatal(err)
			}
		}
		afterFlow, err = flow.SetMerger(firstFlow, tempFlow)
		if err != nil {
			log.Fatal(err)
		}
	} else {
		afterFlow = firstFlow
	}
	secondFlow, err := flow.SetPartitioner(afterFlow, 150, 150)
	if err != nil {
		log.Fatal(err)
	}

	// Send both flows each one to one port. Queues will be added automatically.
	err = flow.SetSender(afterFlow, uint8(outport1))
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

func heavyFunc(currentPacket *packet.Packet, context flow.UserContext) {
	currentPacket.ParseL3()
	ipv4 := currentPacket.GetIPv4()
	if ipv4 != nil {
		T := (ipv4.DstAddr)
		for j := uint(0); j < load; j++ {
			T += uint32(j)
		}
		ipv4.SrcAddr = 263 + (T)
	}
}
