// Copyright 2017 Intel Corporation.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/intel-go/yanff/flow"
	"github.com/intel-go/yanff/packet"
)

var (
	outport uint
	inport  uint

	cloneNumber uint
)

func main() {
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
	f1, err := flow.SetReceiver(uint8(inport))
	if err != nil {
		log.Fatal(err)
	}

	var pdp pcapdumperParameters
	err = flow.SetHandler(f1, pcapdumper, &pdp)
	if err != nil {
		log.Fatal(err)
	}

	// Send packets to control speed. One queue will be added automatically.
	err = flow.SetSender(f1, uint8(outport))
	if err != nil {
		log.Fatal(err)
	}

	err = flow.SystemStart()
	if err != nil {
		log.Fatal(err)
	}
}

type pcapdumperParameters struct {
	f *os.File
}

func (pd pcapdumperParameters) Copy() interface{} {
	filename := fmt.Sprintf("dumped%d.pcap", cloneNumber)
	f, err := os.Create(filename)
	if err != nil {
		fmt.Println("Cannot create file: ", err)
		os.Exit(0)
	}
	cloneNumber++
	err = packet.WritePcapGlobalHdr(f)
	if err != nil {
		log.Fatal(err)
	}
	pdp := pcapdumperParameters{f: f}
	return pdp
}

func (pd pcapdumperParameters) Delete() {
	pd.f.Close()
}

func pcapdumper(currentPacket *packet.Packet, context flow.UserContext) {
	pd := context.(pcapdumperParameters)
	err := currentPacket.WritePcapOnePacket(pd.f)
	if err != nil {
		log.Fatal(err)
	}
}
