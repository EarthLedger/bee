package main

import (
	"encoding/json"
	"fmt"
	"github.com/ethersphere/bee/pkg/swarm"
	"io/ioutil"
	"os"
	"sync"
)

var (
	ThreadNumber = 10
)

type Nodes struct {
	Address []string
}

type TopoProximity struct {
	address   string
	proximity [32]int
}

func (t *TopoProximity) toString() string {
	var info string
	info += fmt.Sprintf("address:%s ", t.address)
	for i, count := range t.proximity {
		info += fmt.Sprintf("%d:%d,", i, count)
	}
	return info
}

func main() {
	fmt.Println("topo")
	if len(os.Args) <= 1 {
		fmt.Println("please input <file name>")
		os.Exit(-1)
	}

	filename := os.Args[1]
	jsonFile, err := os.Open(filename)
	if err != nil {
		fmt.Println("open ", filename, " error: ", err)
		os.Exit(1)
	}

	buffer, _ := ioutil.ReadAll(jsonFile)
	var nodes Nodes
	err = json.Unmarshal(buffer, &nodes)
	if err != nil {
		fmt.Println("load ", filename, " error: ", err)
		os.Exit(1)
	}

	topo := calcAllProximity(&nodes)
	if topo != nil {
		printProximity(topo)
	}

}

func calcAllProximity(nodes *Nodes) []TopoProximity {
	var topo = make([]TopoProximity, len(nodes.Address))

	var calcChan = make(chan int, ThreadNumber)
	wg := sync.WaitGroup{}

	go func() {
		for i := 0; i < len(nodes.Address); i++ {
			calcChan <- i
		}
	}()

	wg.Add(len(nodes.Address))

	for j := 0; j < ThreadNumber; j++ {
		go func() {
			for {
				select {
				case index := <-calcChan:
					base := nodes.Address[index]
					sa, err := swarm.ParseHexAddress(base)
					if err != nil {
						println("wrong address at: ", index, "address: ", base)
						os.Exit(-1)
					}
					topo[index].address = base

					for _, neighbor := range nodes.Address {
						dis, err := calcProximity(sa, neighbor)
						if err != nil {
							os.Exit(-1)
						}
						topo[index].proximity[dis] ++
					}

					wg.Done()
				}
			}
		}()
	}

	wg.Wait()

	return topo
}

func printProximity(topo []TopoProximity) {
	for _, node := range topo {
		fmt.Println(node.toString())
	}
}

func calcProximity(base swarm.Address, neighbor string) (uint8, error) {
	address, err := swarm.ParseHexAddress(neighbor)
	if err != nil {
		println("parse address error ", err)
		return 0, err
	}

	return swarm.Proximity(base.Bytes(), address.Bytes()), nil
}
