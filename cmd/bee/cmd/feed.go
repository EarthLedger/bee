// Copyright 2020 The Swarm Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package cmd

import (
	"math/big"
	"math/rand"

	"github.com/spf13/cobra"

	"github.com/ethersphere/bee/pkg/cac"
	"github.com/ethersphere/bee/pkg/swarm"
)

// input includes: nodes address list, a 4k chunk file, we generate chunk address
// return if this chunk address is closed with any nodes
// close judge is PO 20
func (c *command) initFeedCmd() {
	v := &cobra.Command{
		Use:   "feed",
		Short: "Feed specified nodes",
		Run: func(cmd *cobra.Command, args []string) {
			var closetChunk swarm.Chunk
			var closetPo uint8
			var matchIdx int

			cmd.Println("feed cmd")

			closetDistance := big.NewInt(0)

			addr := args[0]
			cmd.Println(addr)
			node := swarm.MustParseHexAddress(addr)

			for i := 0; i <= 1000000; i++ {
				//cmd.Println(i)
				ch := GenerateRandomChunk()
				//cmd.Println(ch)
				dist, _ := swarm.Distance(node.Bytes(), ch.Address().Bytes())
				//cmd.Println(dist)
				if len(closetDistance.Bits()) == 0 || dist.Cmp(closetDistance) == -1 {
					closetDistance = dist
					matchIdx = i
					closetChunk = ch
					closetPo = swarm.Proximity(node.Bytes(), ch.Address().Bytes())
					cmd.Println("found closer", matchIdx, closetPo)
				}
			}

			cmd.Println("search done!")
			closetPo = swarm.Proximity(node.Bytes(), closetChunk.Address().Bytes())
			cmd.Println("closet chunk:", closetDistance, closetPo, closetChunk)

			// TODO: store this chunk for node

		},
	}
	v.SetOut(c.root.OutOrStdout())
	c.root.AddCommand(v)
}

func GenerateRandomChunk() swarm.Chunk {
	data := make([]byte, swarm.ChunkSize)
	_, _ = rand.Read(data)
	ch, _ := cac.New(data)
	return ch
}
