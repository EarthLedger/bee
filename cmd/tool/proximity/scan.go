package main

import (
	"fmt"
	"github.com/ethersphere/bee/cmd/tool/common"
)

type Neighbor struct {
	Name    string
	Address string
	Success bool
	Po      uint8
}

func (n Neighbor) allToString() string {
	return fmt.Sprintf("name: %s, address: %s, po: %d, success: %t", n.Name, n.Address, n.Po, n.Success)
}

func (n Neighbor) addressToString() string {
	return fmt.Sprintf("%s", n.Address)
}

func load_all_address(filename string) []Neighbor {
	bins:=common.LoadFile(filename)
	if len(bins.Address)<=0{
		return nil
	}

	ns := make([]Neighbor,len(bins.Address))
	for i,addr:=range bins.Address{
		ns[i].Address = addr
	}

	return ns
}
