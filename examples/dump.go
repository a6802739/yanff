// Copyright 2017 Intel Corporation.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/intel-go/yanff/flow"
	"github.com/intel-go/yanff/packet"
)

var (
	outport uint
	inport  uint
)

func main() {
	dumptype := flag.Uint("dumptype", 0, "dumping format type (0 - dumper function, 1 - hex, 2 - pcap file)")
	flag.UintVar(&outport, "outport", 1, "port for sender")
	flag.UintVar(&inport, "inport", 0, "port for receiver")
	flag.Parse()

	// Initialize YANFF library at 10 available cores
	config := flow.Config{
		CPUList: "0-9",
	}
	err := flow.SystemInit(&config)
	if err != nil {
		log.Fatal(err)
	}

	// Receive packets from zero port. One queue will be added automatically.
	firstFlow, err := flow.SetReceiver(uint8(inport))
	if err != nil {
		log.Fatal(err)
	}

	// Separate each 50000000th packet for dumping
	secondFlow, err := flow.SetPartitioner(firstFlow, 50000000, 1)
	if err != nil {
		log.Fatal(err)
	}

	// Dump separated packet. By default function dumper() is used.
	switch *dumptype {
	case 1:
		err = flow.SetHandler(secondFlow, hexdumper, nil)
		if err != nil {
			log.Fatal(err)
		}
	case 2:
		// Writer closes flow
		err = flow.SetWriter(secondFlow, "out.pcap")
		if err != nil {
			log.Fatal(err)
		}
	default:
		err = flow.SetHandler(secondFlow, dumper, nil)
		if err != nil {
			log.Fatal(err)
		}
	}

	// All cases except SetWriter require to merge partitioned packets to original flow
	var output *flow.Flow
	if *dumptype == 2 {
		output = firstFlow
	} else {
		output, err = flow.SetMerger(firstFlow, secondFlow)
		if err != nil {
			log.Fatal(err)
		}
	}
	err = flow.SetSender(output, uint8(outport))
	if err != nil {
		log.Fatal(err)
	}

	err = flow.SystemStart()
	if err != nil {
		log.Fatal(err)
	}
}

func dumper(currentPacket *packet.Packet, context flow.UserContext) {
	var tcp *packet.TCPHdr
	var udp *packet.UDPHdr
	var icmp *packet.ICMPHdr

	fmt.Printf("%v", currentPacket.Ether)
	ipv4, ipv6, arp := currentPacket.ParseAllKnownL3()
	if ipv4 != nil {
		fmt.Printf("%v", ipv4)
		tcp, udp, icmp = currentPacket.ParseAllKnownL4ForIPv4()
	} else if ipv6 != nil {
		fmt.Printf("%v", ipv6)
		tcp, udp, icmp = currentPacket.ParseAllKnownL4ForIPv6()
	} else if arp != nil {
		fmt.Printf("%v", arp)
	} else {
		fmt.Println("    Unknown L3 protocol")
	}

	if tcp != nil {
		fmt.Printf("%v", tcp)
	} else if udp != nil {
		fmt.Printf("%v", udp)
	} else if icmp != nil {
		fmt.Printf("%v", icmp)
	} else {
		fmt.Println("        Unknown L4 protocol")
	}
	fmt.Println("----------------------------------------------------------")
}

func hexdumper(currentPacket *packet.Packet, context flow.UserContext) {
	fmt.Printf("Raw bytes=%x\n", currentPacket.GetRawPacketBytes())
}
