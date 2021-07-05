package main

import (
	"github.com/ethersphere/bee/pkg/swarm"
	"github.com/ethersphere/bee/pkg/topology/pslice"
	"testing"
)

var (
	BaseAddress = "5a3dad39940b63b48f585e3d95da5a42fac00f62d03b791f621376d9d963eaa8"
	filename    = "ppp.log"
)

func Test_deeps(t *testing.T) {
	address, err := swarm.ParseHexAddress(BaseAddress)
	if err != nil {
		t.Error("parse address error")
	}

	peers := pslice.New(int(swarm.MaxBins), address)
	ns := load_all_address(filename)
	for _, n := range ns {
		address, err := swarm.ParseHexAddress(n.Address)
		if err != nil {
			t.Error("parse address error")
		}
		peers.Add(address)
	}
}
