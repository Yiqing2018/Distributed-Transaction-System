package main

import (
	"fmt"
	"strconv"
	"strings"
	"sync"
)

var circleNode []Node

type Node struct {
	value string
}

type Graph struct {
	edges map[Node][]*Node
	lock  sync.RWMutex
}

func (graph *Graph) addEdge(start, end *Node) {
	graph.lock.Lock()
	defer graph.lock.Unlock()

	if graph.edges == nil {
		graph.edges = make(map[Node][]*Node)
	}

	// start -> end
	for _, item := range graph.edges[*start] {
		if item.value == end.value {
			return
		}
	}

	graph.edges[*start] = append(graph.edges[*start], end)
}

func findSmallestTID() string {
	result := 9999
	for _, item := range circleNode {
		itemInt, _ := strconv.Atoi(item.value)
		if itemInt < result {
			result = itemInt
		}
	}

	return strconv.Itoa(result)
}

func (graph *Graph) handleDelete(deleteTID string) {
	// 1.1 delete smallestTID when smallestTID is Boss
	delete(graph.edges, Node{deleteTID})

	// 1.2 delete smallestTID when smallestTID is Didi
	for eachNode := range graph.edges {
		neighbors := graph.edges[eachNode]
		var newSlice []*Node
		for _, neighbor := range neighbors {
			if neighbor.value == deleteTID {
				continue
			}
			newSlice = append(newSlice, neighbor)
		}
		if len(newSlice) == 0 {
			delete(graph.edges, eachNode)
		} else {
			graph.edges[eachNode] = newSlice
		}
	}
}

func (graph *Graph) handleCircle() {
	graph.lock.Lock()
	defer graph.lock.Unlock()

	// 1 send to server to abort smallestTID
	smallestTID := findSmallestTID()
	graph.handleDelete(smallestTID)
	sendServer("doAbort" + " " + smallestTID)
}

func (graph *Graph) traversalGraph() {
	graph.lock.RLock()
	defer graph.lock.RUnlock()
	str := ""

	for eachNode := range graph.edges {
		if eachNode.value == "-1" || eachNode.value == "-2" {
			continue
		}
		str += eachNode.value + " -> "
		neighbors := graph.edges[eachNode]

		for _, neighbor := range neighbors {
			str += neighbor.value + " "
		}
		str += "\n"
	}

	fmt.Println(str)
}

func hasValue(node *Node, stack []Node) bool {
	for _, a := range stack {
		if a == *node {
			return true
		}
	}
	return false
}

func findIndex(node *Node, stack []Node) int {
	for i, a := range stack {
		if a == *node {
			return i
		}
	}
	return len(stack)
}

func dfs(node *Node, graph Graph, isVisit map[Node]bool, stack []Node) {
	isVisit[*node] = true
	stack = append(stack, *node)

	for _, neighbor := range graph.edges[*node] {
		if !hasValue(neighbor, stack) {
			if !isVisit[*neighbor] {
				if _, ok := graph.edges[*neighbor]; ok {
					dfs(neighbor, graph, isVisit, stack)
				}
			}
		} else {
			if _, ok := graph.edges[*neighbor]; ok {
				index := findIndex(neighbor, stack)
				for _, item := range stack[index:] {
					circleNode = append(circleNode, item)
				}
			}
		}
	}

	if len(stack) > 0 {
		stack = stack[:len(stack)-1]
	}

}

func (testGraph *Graph) addTestEdge(testEdge string) {
	startNode := Node{strings.Split(testEdge, " ")[0]}
	endNode := Node{strings.Split(testEdge, " ")[1]}

	testGraph.addEdge(&startNode, &endNode)
}

func (testGraph *Graph) hasCircle() bool {
	isVisit := make(map[Node]bool)
	for eachNode := range testGraph.edges {
		isVisit[eachNode] = false
	}

	var stack []Node
	for eachNode := range isVisit {
		if !isVisit[eachNode] {
			dfs(&eachNode, *testGraph, isVisit, stack)
		}
	}

	if len(circleNode) != 0 {
		return true
	}

	return false
}

func (testGraph *Graph) checkCircle() {
	for testGraph.hasCircle() {
		fmt.Println("Deadlock Detected:", circleNode)
		fmt.Println("The Graph (With DeadLock) is:")
		testGraph.traversalGraph()
		testGraph.handleCircle()
		var tempCircleNode []Node
		circleNode = tempCircleNode
	}

	fmt.Println("The Graph (Without DeadLock) is:")
	testGraph.traversalGraph()
}

func (testGraph *Graph) initGraph() {
	initLeft := Node{"-1"}
	initRight := Node{"-2"}

	testGraph.addEdge(&initLeft, &initRight)
}
