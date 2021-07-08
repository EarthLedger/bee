package common

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)


type Nodes struct {
	Address []string `json:"address"`
}

func LoadFile(filename string) *Nodes {
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
	return &nodes
}