package graph

import (
	"errors"
	"fmt"
	"log"
	"math"
	"math/rand"
	"time"

	"golang.org/x/exp/slices"
)

type dv map[string]float64

type broadcast struct {
	from *node
	dv   dv
}

type node struct {
	name  string
	conns map[*node]float64

	// distance vectors
	dv dv

	changeChan chan broadcast
}

func NewNode(name string) *node {
	changeChan := make(chan broadcast, 100)
	n := &node{
		name:       name,
		conns:      make(map[*node]float64),
		dv:         make(map[string]float64),
		changeChan: changeChan,
	}

	go listenForChange(n)

	return n
}

func listenForChange(n *node) {
	for {
		log.Printf("%s is listening for change...\n", n.name)

		broadcast := <-n.changeChan
		log.Printf("%s got broadcast from %s \n", n.name, broadcast.from.name)
		var changed bool
		for bNode, bCost := range broadcast.dv {
			if bNode == n.name {
				continue
			}
			oldCost, found := n.dv[bNode]
			if !found {
				oldCost = math.Inf(1)
			}

			newCost := bCost + n.conns[broadcast.from]
			if newCost < oldCost {
				n.dv[bNode] = newCost
				changed = true
			}
		}

		// mimic network delay
		dt := rand.Float32() * 3
		time.Sleep(time.Duration(dt) * time.Second)
		log.Printf("%s updated its DV in %f seconds! \n", n.name, dt)

		if changed {
			log.Printf("%s dv was changed, broadcasting to neighbors... \n", n.name)
			n.broadcast()
		}
	}
}

func (n *node) GetNeighbor() []*node {
	var neighbors []*node
	for conn := range n.conns {
		neighbors = append(neighbors, conn)
	}

	return neighbors
}

func (n *node) GetNeighborsNames() []string {
	var neighbors []string
	for conn := range n.conns {
		neighbors = append(neighbors, conn.name)
	}

	return neighbors
}

func (n *node) addEdge(to *node, cost float64) {
	n.conns[to] = cost
	n.dv[to.name] = cost
}

func (n *node) broadcast() {
	for conn := range n.conns {
		log.Printf("%s sent its dv to %s \n", n.name, conn.name)
		conn.changeChan <- broadcast{from: n, dv: n.dv}
	}
}

type graph struct {
	nodes map[string]*node
}

func NewGraph() *graph {
	return &graph{
		nodes: make(map[string]*node),
	}
}

func (g *graph) AddNode(n *node) error {
	_, found := g.nodes[n.name]
	if found {
		return errors.New("nodes with duplicate names are now allowed")
	}

	g.nodes[n.name] = n

	return nil
}

func (g *graph) AddEdge(from, to string, cost float64) error {
	if from == to {
		return errors.New("nodes can not connect to themselves")
	}
	fromNode, found := g.nodes[from]
	if !found {
		return fmt.Errorf("node %s not available in the graph", from)
	} else if slices.Contains(fromNode.GetNeighborsNames(), to) {
		return fmt.Errorf("edge between %s and %s already exists", from, to)
	}

	toNode, found := g.nodes[to]
	if !found {
		return fmt.Errorf("node %s not available in the graph", to)
	} else if slices.Contains(toNode.GetNeighborsNames(), from) {
		return fmt.Errorf("edge between %s and %s already exists", to, from)
	}

	fromNode.addEdge(toNode, cost)
	toNode.addEdge(fromNode, cost)

	return nil
}

func (g *graph) Broadcast() {
	for _, n := range g.nodes {
		n.broadcast()
	}
}

func (g *graph) Describe() string {
	desc := "Graph: \n"
	for _, n := range g.nodes {
		neighbors := []string{}
		for to, cost := range n.dv {
			neighbors = append(neighbors, fmt.Sprintf("%s:%.1f", to, cost))
		}
		desc += fmt.Sprintf("Node: %s, Neighbors: %s \n", n.name, neighbors)
	}

	return desc
}
