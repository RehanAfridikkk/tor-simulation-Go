// main.go

package main

import (
	"fmt"
	"math/rand"
	"time"
	"tor/structure"
)

func main() {
	rand.Seed(time.Now().UnixNano())

	var nodeCount int
	fmt.Print("Enter the number of nodes: ")
	fmt.Scan(&nodeCount)

	network := structure.NewNetwork(nodeCount)

	srcNodeID := rand.Intn(nodeCount) + 1
	destNodeID := rand.Intn(nodeCount) + 1

	msg := "Hello, anonymous world!"

	network.SendMessage(srcNodeID, destNodeID, msg)

	time.Sleep(time.Second)
}
