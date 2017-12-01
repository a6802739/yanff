// Call insmod ./dpdk/dpdk-17.08/x86_64-native-linuxapp-gcc/kmod/rte_kni.ko lo_mode=lo_mode_fifo_skb
// before this. It will make a loop of packets inside KNI device and "send" will send received packets.
// Other variants of rte_kni.ko configuration can be found here:
// http://dpdk.org/doc/guides/sample_app_ug/kernel_nic_interface.html

// Need to call "ifconfig myKNI 111.111.11.11" while running this example to allow other applications
// to receive packets from "111.111.11.11" address

package main

import (
	"log"

	"github.com/intel-go/yanff/flow"
)

func main() {
	config := flow.Config{
		// Is required for KNI
		NeedKNI: true,
		CPUList: "0-7",
	}

	err := flow.SystemInit(&config)
	if err != nil {
		log.Fatal(err)
	}
	// (port of device, core (not from YANFF set) which will handle device, name of device)
	kni := flow.CreateKniDevice(1, 20, "myKNI")

	fromEthFlow, err := flow.SetReceiver(0)
	if err != nil {
		log.Fatal(err)
	}
	err = flow.SetSender(fromEthFlow, kni)
	if err != nil {
		log.Fatal(err)
	}

	fromKNIFlow, err := flow.SetReceiver(kni)
	if err != nil {
		log.Fatal(err)
	}
	err = flow.SetSender(fromKNIFlow, 1)
	if err != nil {
		log.Fatal(err)
	}

	err = flow.SystemStart()
	if err != nil {
		log.Fatal(err)
	}
}
