package main

import (
	"fmt"
	"github.com/ethersphere/bee/cmd/tool/common"
	"github.com/ethersphere/bee/pkg/swarm"
	"os"
	"strconv"
	"sync"
)

var (
	ThreadNumber = 10
)

type TopoProximity struct {
	index int
	address   string
	proximity [32]int
}

func (t *TopoProximity) headString() string {
	var info string
	info += fmt.Sprintf("%-8s", "index")
	info += fmt.Sprintf("%-68s", "address")
	for i, _ := range t.proximity {
		info += fmt.Sprintf("deep%-8d", i)
	}
	return info
}


func (t *TopoProximity) toString() string {
	var info string
	info += fmt.Sprintf("%-8d", t.index)
	info += fmt.Sprintf("%-68s", t.address)
	for _, count := range t.proximity {
		info += fmt.Sprintf("%-12d", count)
	}
	return info
}

func main() {
	if len(os.Args) != 4 && len(os.Args) != 6 {
		fmt.Println("please input <op> " +
			"case: topo <file name> <thread number>" +
			"case: find <file name> <thread number> <address> <proximity>")
		os.Exit(-1)
	}

	op := os.Args[1]
	filename := os.Args[2]
	threads, err := strconv.Atoi(os.Args[3])
	if err != nil {
		fmt.Println("wrong thread number, will use default thread number of 10")
	} else {
		ThreadNumber = threads
	}

	bins := common.LoadFile(filename)

	if op == "topo" {
		topo(bins)
	} else if op == "find" {
		base := os.Args[4]
		proximity, err := strconv.Atoi(os.Args[5])
		if err != nil {
			fmt.Println("wrong proximity")
			os.Exit(-1)
		}
		find(base, bins, proximity)
	} else {
		fmt.Println("please input <op> " +
			"case: topo <file name> <thread number>" +
			"case: find <file name> <address> <proximity>")
		os.Exit(-1)
	}
}


func topo(bins *common.Nodes) {
	topo := calcAllProximity(bins)
	if topo != nil {
		printProximity(topo)
	}
}

func find(base string, bins *common.Nodes, proximity int)  {
	nodesArray := findIn(base,bins,proximity)
	for i,arr:=range nodesArray{
		info:=fmt.Sprintf("index:%d, total:%d ",i,len(arr.Address))
		for _,addr:=range arr.Address{
			info+=fmt.Sprintf(" %s ",addr)
		}
		fmt.Println(info)
	}
}

func findIn(base string, bins *common.Nodes, proximity int) []common.Nodes {
	var nodesArray []common.Nodes
	ones := findOne(base, bins, proximity)
	if len(ones.Address) == 1 {
		var nodes common.Nodes
		nodes.Address = append(nodes.Address, ones.Address[0])
		nodes.Address = append(nodes.Address, base)
		nodesArray = append(nodesArray, nodes)
		return nodesArray
	} else if len(ones.Address) == 0 {
		var nodes common.Nodes
		nodes.Address = append(nodes.Address, base)
		nodesArray = append(nodesArray, nodes)
		return nodesArray
	}

	for _, elem := range ones.Address {
		middle := findIn(elem, ones, proximity)
		if len(middle) == 1{
			var node common.Nodes
			node.Address = append(node.Address, middle[0].Address...)
			node.Address = append(node.Address,base)
			nodesArray = append(nodesArray, node)
			continue
		}

		for _,mid:=range middle{
			var node common.Nodes
			node.Address = append(node.Address, mid.Address...)
			node.Address = append(node.Address,base)
			nodesArray = append(nodesArray, node)
		}
	}

	return nodesArray
}

func findOne(base string, bins *common.Nodes, proximity int) *common.Nodes {
	sa, err := swarm.ParseHexAddress(base)
	if err != nil {
		println("wrong base address ", base)
		os.Exit(-1)
	}

	var finds common.Nodes
	for _, neighbor := range bins.Address {
		dis, err := calcProximity(sa, neighbor)
		if err != nil {
			os.Exit(-1)
		}

		if proximity <= int(dis) {
			finds.Address = append(finds.Address, neighbor)
		}
	}

	return &finds
}

func calcAllProximity(nodes *common.Nodes) []TopoProximity {
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
					topo[index].index = index
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
	var temp TopoProximity
	fmt.Println(temp.headString())
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
