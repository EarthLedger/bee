// Copyright 2020 The Swarm Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package cmd

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/big"
	"math/rand"
	"os"
	"strconv"

	"github.com/spf13/cobra"

	"github.com/ethersphere/bee/pkg/cac"
	"github.com/ethersphere/bee/pkg/swarm"
)

type ClosetChunk struct {
	node     swarm.Address
	chunk    swarm.Chunk
	count    int
	distance *big.Int
	po       uint8
}

// input includes: nodes address list, a 4k chunk file, we generate chunk address
// return if this chunk address is closed with any nodes
// close judge is PO 20
func (c *command) initFeedCmd() {
	v := &cobra.Command{
		Use:   "feed",
		Short: "feed [address json file] [try count]",
		Run: func(cmd *cobra.Command, args []string) {

			if len(args) != 2 {
				cmd.Println("error, miss params")
				return
			}
			addrFile := args[0]
			nodeMap := parseAddressJsonFile(addrFile)
			count, _ := strconv.Atoi(args[1])

			findClosetChunk(cmd, nodeMap, count)

			for k, data := range nodeMap {
				filePath := fmt.Sprintf("./chunks/%s", k)
				os.MkdirAll(filePath, os.ModePerm)
				fileName := fmt.Sprintf("./chunks/%s/%d-%s", k, data.po, data.chunk.Address())
				err := ioutil.WriteFile(fileName, data.chunk.Data(), 0644)
				check(err)
			}
		},
	}
	v.SetOut(c.root.OutOrStdout())
	c.root.AddCommand(v)
}

func parseAddressJsonFile(filePath string) map[string]*ClosetChunk {
	var nodeMap map[string]*ClosetChunk = make(map[string]*ClosetChunk)
	addrFile, e := ioutil.ReadFile(filePath)
	check(e)
	var data []string
	err := json.Unmarshal(addrFile, &data)
	check(err)

	for _, addr := range data {
		node := swarm.MustParseHexAddress(addr)
		chunk := GenerateRandomChunk()
		distance, _ := swarm.Distance(node.Bytes(), chunk.Address().Bytes())
		nodeMap[addr] = &ClosetChunk{
			node:     node,
			chunk:    chunk,
			count:    0,
			distance: distance,
			po:       swarm.Proximity(node.Bytes(), chunk.Address().Bytes()),
		}
	}

	return nodeMap
}

func GenerateRandomChunk() swarm.Chunk {
	data := make([]byte, swarm.ChunkSize)
	_, _ = rand.Read(data)
	ch, _ := cac.New(data)
	return ch
}

func findClosetChunk(cmd *cobra.Command, nodeMap map[string]*ClosetChunk, count int) {
	for i := 0; i <= count; i++ {
		//cmd.Printf("\r%d ", i)
		ch := GenerateRandomChunk()
		for k, data := range nodeMap {
			po := swarm.Proximity(data.node.Bytes(), ch.Address().Bytes())
			if po > data.po {
				data.chunk = ch
				data.count = i
				data.po = po
				cmd.Println("found closer", i, data.po, k)
			}
		}
	}

	cmd.Println("search done!")
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}
