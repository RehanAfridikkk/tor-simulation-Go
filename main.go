// main.go

package main

import (
	"math/rand"
	"time"
	"tor/structure"
)

func main() {
	rand.Seed(time.Now().UnixNano())

	nodeCount := 5
	network := structure.NewNetwork(nodeCount)

	srcNodeID := rand.Intn(nodeCount) + 1
	destNodeID := rand.Intn(nodeCount) + 1

	msg := "Hello, anonymous world!"

	network.SendMessage(srcNodeID, destNodeID, msg)

	time.Sleep(time.Second)
}
