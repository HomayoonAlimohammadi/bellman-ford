package main

import (
	"fmt"
	"log"
	"time"

	"go-test/bellman-ford/graph"
)

func main() {
	g := graph.NewGraph()
	check(g.AddNode(graph.NewNode("A")))
	check(g.AddNode(graph.NewNode("B")))
	check(g.AddNode(graph.NewNode("C")))
	check(g.AddNode(graph.NewNode("D")))
	check(g.AddNode(graph.NewNode("E")))

	check(g.AddEdge("A", "C", 1))
	check(g.AddEdge("B", "C", 4))
	check(g.AddEdge("C", "E", 1))
	check(g.AddEdge("B", "E", 1))
	check(g.AddEdge("C", "D", 3))
	check(g.AddEdge("D", "E", 1))

	time.Sleep(2 * time.Second)
	fmt.Println()
	fmt.Println(g.Describe())

	g.Broadcast()

	time.Sleep(10 * time.Second)
	fmt.Println()
	fmt.Println(g.Describe())
}

func check(err error) {
	if err != nil {
		log.Println(err)
	}
}
