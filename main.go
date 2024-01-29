package main

import (
	"fmt"
	"math/rand"
	"net/http"
	"time"
)

type Node struct {
	ID       int
	Inbound  chan string
	Outbound chan string
}

type Network struct {
	Nodes []*Node
}

func NewNode(id int) *Node {
	return &Node{
		ID:       id,
		Inbound:  make(chan string),
		Outbound: make(chan string),
	}
}

func NewNetwork(nodeCount int) *Network {
	network := &Network{}
	for i := 1; i <= nodeCount; i++ {
		node := NewNode(i)
		network.Nodes = append(network.Nodes, node)
		go network.startNode(node)
		go network.startServer(node)
	}
	return network
}

func (n *Network) startNode(node *Node) {
	for {
		select {
		case msg := <-node.Inbound:
			fmt.Printf("Node %d received message: %s\n", node.ID, msg)
			// Simulate message forwarding (pseudo-anonymity)
			time.Sleep(time.Millisecond * time.Duration(rand.Intn(100)))
			// Forward the message to a random node
			destNodeID := rand.Intn(len(n.Nodes)) + 1
			n.Nodes[destNodeID-1].Inbound <- msg
		case msg := <-node.Outbound:
			// Simulate message delivery
			fmt.Printf("Node %d delivered message: %s\n", node.ID, msg)
		}
	}
}

func (n *Network) startServer(node *Node) {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		msg := r.URL.Path[1:]
		node.Inbound <- msg
		fmt.Fprintf(w, "Node %d received message: %s", node.ID, msg)
	})

	port := fmt.Sprintf(":%d", 9000+node.ID)
	fmt.Printf("Node %d listening on port %s\n", node.ID, port)
	http.ListenAndServe(port, nil)
}

func (n *Network) sendMessage(srcNodeID int, destNodeID int, msg string) {
	n.Nodes[srcNodeID-1].Outbound <- msg
	n.Nodes[destNodeID-1].Inbound <- msg
}

func main() {
	rand.Seed(time.Now().UnixNano())

	nodeCount := 5
	network := NewNetwork(nodeCount)

	srcNodeID := rand.Intn(nodeCount) + 1
	destNodeID := rand.Intn(nodeCount) + 1

	msg := "Hello, anonymous world!"

	network.sendMessage(srcNodeID, destNodeID, msg)

	time.Sleep(time.Second)
}
