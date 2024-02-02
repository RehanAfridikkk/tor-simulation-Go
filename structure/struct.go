// structure/struct.go

package structure

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
	Nodes   []*Node
	Counter int
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
		go network.StartNode(node)
		go network.StartServer(node)
	}
	return network
}

func (n *Network) StartNode(node *Node) {
	for {
		select {
		case msg := <-node.Inbound:
			fmt.Printf("Node %d received message: %s\n", node.ID, msg)

			// Check if the current node is the action node (node 4)
			if node.ID == 4 {
				fmt.Println("Action node reached. Making HTTP request to google.com.")

				// Make an HTTP request to google.com
				resp, err := http.Get("http://google.com")
				if err != nil {
					fmt.Println("Error making HTTP request:", err)
				} else {
					defer resp.Body.Close()
					fmt.Println("HTTP request to google.com successful.")
				}

				fmt.Println("Stopping.")
				return
			}

			// Simulate message forwarding (pseudo-anonymity)
			time.Sleep(time.Millisecond * time.Duration(rand.Intn(100)))
			// Forward the message to the next node in a deterministic sequence
			destNodeID := (node.ID % len(n.Nodes)) + 1
			n.Nodes[destNodeID-1].Inbound <- msg

			// Increment the counter
			n.Counter++

			// Check if the message has circulated through all nodes twice
			if n.Counter == 2*len(n.Nodes) {
				fmt.Println("Message circulated through all nodes twice. Stopping.")
				return
			}

		case msg := <-node.Outbound:
			// Simulate message delivery
			fmt.Printf("Node %d delivered message: %s\n", node.ID, msg)
		}
	}
}
func (n *Network) StartServer(node *Node) {
	path := fmt.Sprintf("/node%d", node.ID)

	http.HandleFunc(path, func(w http.ResponseWriter, r *http.Request) {
		msg := r.URL.Path[len(path):]
		node.Inbound <- msg
		fmt.Fprintf(w, "Node %d received message: %s", node.ID, msg)
	})

	port := fmt.Sprintf(":%d", 9000+node.ID)
	fmt.Printf("Node %d listening on path %s\n", node.ID, path)
	http.ListenAndServe(port, nil)
}

func (n *Network) SendMessage(srcNodeID int, destNodeID int, msg string) {
	n.Nodes[srcNodeID-1].Outbound <- msg
	n.Nodes[destNodeID-1].Inbound <- msg
}
