// Copyright 2017 Intel Corporation.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"flag"
	"log"

	"github.com/intel-go/yanff/flow"
)

var inport, outport uint
var cores string

// This is a test for pure send/receive performance measurements. No
// other functions used here.
func main() {
	flag.UintVar(&inport, "inport", 0, "Input port number")
	flag.UintVar(&outport, "outport", 0, "Output port number")
	flag.StringVar(&cores, "cores", "0-15", "Specifies CPU cores to be used by YANFF library")
	flag.Parse()

	// Initialize YANFF library to use specified number of CPU cores
	config := flow.Config{
		CPUList: cores,
	}
	err := flow.SystemInit(&config)
	if err != nil {
		log.Fatal(err)
	}

	// Receive packets from input port. One queue will be added automatically.
	f, err := flow.SetReceiver(uint8(inport))
	if err != nil {
		log.Fatal(err)
	}

	err = flow.SetSender(f, uint8(outport))
	if err != nil {
		log.Fatal(err)
	}

	err = flow.SystemStart()
	if err != nil {
		log.Fatal(err)
	}
}
