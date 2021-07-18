package main

import (
	"encoding/json"
	"fmt"
	"github.com/ethersphere/bee/cmd/tool/common"
	"github.com/ethersphere/bee/pkg/swarm"
	"io/ioutil"
	"os"
	"strconv"
)

func main() {
	if len(os.Args) <= 2 {
		fmt.Println("please input <base node swarm address> <peer swarm address> ")
		os.Exit(-1)
	}

	base := os.Args[1]
	peerAddress := os.Args[2]
	baseAddress, err := swarm.ParseHexAddress(base)
	if err != nil {
		fmt.Println("input base address parse error: ", err)
		os.Exit(-1)
	}

	po, err := calcProximity(baseAddress, peerAddress)
	if err != nil {
		fmt.Println("input peer address parse error: ", err)
		os.Exit(-1)
	}

	fmt.Printf("proximity: %d\n", po)
}

func main_old() {
	if len(os.Args) <= 3 {
		fmt.Println("please input <base address> <file name> <proximity>")
		os.Exit(-1)
	}

	address := os.Args[1]
	filename := os.Args[2]
	proximity, _ := strconv.Atoi(os.Args[3])
	targetFileName :=  os.Args[4]

	fmt.Println("peer address ", address)
	addr, err := swarm.ParseHexAddress(address)
	if err != nil {
		fmt.Println("input base address parse error: ", err)
		os.Exit(-1)
	}

	calcNeighborProximity(addr, filename, proximity, targetFileName)
}

func calcNeighborProximity(base swarm.Address, filename string, proximity int, targetFileName string) {
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
	saveNeighborByPo(ns, proximity, targetFileName)
}

func calcProximity(base swarm.Address, neighbor string) (uint8, error) {
	address, err := swarm.ParseHexAddress(neighbor)
	if err != nil {
		println("parse address error ", err)
		return 0, err
	}

	return swarm.Proximity(base.Bytes(), address.Bytes()), nil
}

func saveNeighborByPo(ns []Neighbor, proximity int, targetFileName string) {
	if proximity == -1 {
		saveAllNeighborInfo(ns, targetFileName)
	} else {
		saveNeighborAddressByPo(ns, proximity, targetFileName)
	}
}

func saveAllNeighborInfo(ns []Neighbor, targetFileName string) {
	var info []string
	for i := 0; i < len(ns); i++ {
		info = append(info,  fmt.Sprint(ns[i].allToString()))
	}

	file, _ := json.MarshalIndent(info, "", " ")
	_ = ioutil.WriteFile(targetFileName, file, 0644)
}

func saveNeighborAddressByPo(ns []Neighbor, proximity int, targetFileName string) {
	var bins common.Nodes
	for i := 0; i < len(ns); i++ {
		if proximity == int(ns[i].Po) {
			//fmt.Println(ns[i].addressToString())
			bins.Address = append(bins.Address, ns[i].Address)
		}
	}

	file, _ := json.MarshalIndent(bins, "", " ")
	_ = ioutil.WriteFile(targetFileName, file, 0644)
}
