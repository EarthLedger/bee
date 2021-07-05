package main

import (
	"fmt"
	"github.com/ethersphere/bee/pkg/swarm"
	"os"
)

func main() {
	fmt.Println("proximity")
	if len(os.Args) <= 3 {
		fmt.Println("please input <base address> <file name> <load or other>")
		os.Exit(-1)
	}

	address := os.Args[1]
	filename := os.Args[2]
	action := os.Args[3]

	fmt.Println("peer address ", address)
	addr, err := swarm.ParseHexAddress(address)
	if err != nil {
		fmt.Println("input base address parse error: ", err)
		os.Exit(-1)
	}

	calcNeighborProximity(addr, filename, action)
}

func calcNeighborProximity(base swarm.Address, filename string, action string) {
	ns := load_all_address(filename)
	for i := range ns {
		po, err := calcProximity(base, ns[i].Address)
		if err != nil {
			ns[i].Success = false
		} else {
			ns[i].Success = true
		}
		ns[i].Po = po
	}
	saveNeighborPo(ns, action)
}

func calcProximity(base swarm.Address, neighbor string) (uint8, error) {
	address, err := swarm.ParseHexAddress(neighbor)
	if err != nil {
		println("parse address error ", err)
		return 0, err
	}

	return swarm.Proximity(base.Bytes(), address.Bytes()), nil
}

func saveNeighborPo(ns []Neighbor, action string) {
	for i := 0; i < len(ns); i++ {
		if action == "load" {
			fmt.Println(ns[i].Address)
		} else {
			fmt.Println(ns[i].tostring())
		}

	}
}
