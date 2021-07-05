package main

import (
	"bufio"
	"fmt"
	"os"
)

type Neighbor struct {
	Name    string
	Address string
	Success bool
	Po      uint8
}

func (n Neighbor) tostring() string {
	return fmt.Sprintf("name: %s, address: %s, po: %d, success: %t", n.Name, n.Address, n.Po, n.Success)
}

func load_all_address(filename string) []Neighbor {
	file, err := os.Open(filename)
	if err != nil {
		fmt.Println("read file error: ", err)
		os.Exit(-1)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)

	ns := []Neighbor{}
	neighbor := Neighbor{}
	bName := true
	for scanner.Scan() {
		if bName {
			text := scanner.Text()
			neighbor.Name = text
			bName = false
		} else {
			text := scanner.Text()
			neighbor.Address = text[12 : 12+40]
			bName = true
			ns = append(ns, neighbor)
		}
	}

	return ns
}
