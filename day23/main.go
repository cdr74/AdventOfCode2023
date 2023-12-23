// ---------------------------------------------------------------------------
// Golang solution for Advent of Code 2023 day 23
// https://adventofcode.com/2023/day/23
//
// This was created using copilot to assist me in learning Go.
//
// Scenario:
//
//	Go on a hike on the map given in put data. Rules
//	1 Slopes can only be passed in direction of the arrow.
//	2 Tiles are ony visited once.
//	3 start top row only . tile
//	4 end bottom row only . tile
//
// Data
//
//	Trails map with tiles, that indicates
//	 (.) path
//	 (#) blocked
//	 (^, >, v, <) slopes, only pass in given direction
//
// Part 1:
//
//	Find the longest hike without loops. Longest path is a NP problem ...
//	test.data -> 94 steps
//
// Part 2:
// ---------------------------------------------------------------------------
package main

import (
	"fmt"
	"time"

	"github.com/cdr74/AdventOfCode2023/utils"
)

// runTest defines whether test.data or actual.data should be used
const runTest bool = false
const TEST_FILE string = "test.data"
const DATA_FILE string = "actual.data"

// -------------------------- Common Code Section ----------------------------

func getNodeID(i, j int) string {
	return fmt.Sprintf("%d_%d", i, j)
}

// Node represents a node in the graph
type Node struct {
	ID        string
	Neighbors []*Node
}

// AddNeighbor adds a neighbor to the node
func (n *Node) AddNeighbor(neighbor *Node) {
	n.Neighbors = append(n.Neighbors, neighbor)
}

// Graph represents the graph
type Graph struct {
	Nodes map[string]*Node
}

// AddNode adds a node to the graph
func (g *Graph) AddNode(node *Node) {
	g.Nodes[node.ID] = node
}

func (g *Graph) print() {
	for _, node := range g.Nodes {
		fmt.Printf("Node: %s\t -> ", node.ID)
		for _, neighbor := range node.Neighbors {
			fmt.Printf("%s ", neighbor.ID)
		}
		fmt.Println()
	}
}

func connectNeighbors(graph *Graph, neighborMap *map[string][]string) {
	for nodeID, neighbors := range *neighborMap {
		for _, neighborID := range neighbors {
			if neighbor, ok := graph.Nodes[neighborID]; ok {
				graph.Nodes[nodeID].AddNeighbor(neighbor)
			}
		}
	}
}

// CreateGraph creates a graph from a 2x2 char array
func CreateGraph(grid []string) *Graph {
	neighborMap := make(map[string][]string)
	graph := &Graph{Nodes: make(map[string]*Node)}
	for i := 0; i < len(grid); i++ {
		for j, c := range grid[i] {
			if c != '#' {
				nodeID := getNodeID(i, j)
				node := &Node{ID: nodeID}
				graph.AddNode(node)

				switch c {
				case '^':
					neighborMap[nodeID] = []string{getNodeID(i-1, j)}
				case '>':
					neighborMap[nodeID] = []string{getNodeID(i, j+1)}
				case 'v':
					neighborMap[nodeID] = []string{getNodeID(i+1, j)}
				case '<':
					neighborMap[nodeID] = []string{getNodeID(i, j-1)}
				default:
					neighborMap[nodeID] = []string{getNodeID(i+1, j), getNodeID(i-1, j), getNodeID(i, j+1), getNodeID(i, j-1)}
				}
			}
		}
	}
	connectNeighbors(graph, &neighborMap)

	return graph
}

func CreateGraph_Part2(grid []string) *Graph {
	neighborMap := make(map[string][]string)
	graph := &Graph{Nodes: make(map[string]*Node)}
	for i := 0; i < len(grid); i++ {
		for j, c := range grid[i] {
			if c != '#' {
				nodeID := getNodeID(i, j)
				node := &Node{ID: nodeID}
				graph.AddNode(node)
				neighborMap[nodeID] = []string{getNodeID(i+1, j), getNodeID(i-1, j), getNodeID(i, j+1), getNodeID(i, j-1)}
			}
		}
	}
	connectNeighbors(graph, &neighborMap)

	return graph
}

// -------------------------- Puzzle part 1 ----------------------------------

// no optimization done, given this is a np problem... not ideal
// memory of visited nodes could help a lot
func dfs(graph *Graph, node *Node, end *Node, visited *map[string]bool, path *[]string, maxLength *int) {
	(*visited)[node.ID] = true
	*path = append(*path, node.ID)

	if node == end {
		if len(*path) > *maxLength {
			*maxLength = len(*path)
			//fmt.Printf("Path: %v - steps %d\n", *path, *maxLength-1)
		}
	}

	for _, neighbor := range node.Neighbors {
		if !(*visited)[neighbor.ID] {
			dfs(graph, neighbor, end, visited, path, maxLength)
		}
	}

	// backtracking
	*path = (*path)[:len(*path)-1]
	(*visited)[node.ID] = false
}

// As longest path is a np hard problem, best bet is depth first search
// with backtracking.
func SolvePart1(g *Graph) int {
	startNodeID := getNodeID(0, 1)
	endNodeID := getNodeID(140, 139)
	if runTest {
		endNodeID = getNodeID(22, 21)
	}
	path := []string{}
	maxLength := 0
	dfs(g, g.Nodes[startNodeID], g.Nodes[endNodeID], &map[string]bool{}, &path, &maxLength)
	return maxLength - 1
}

// -------------------------- Puzzle part 2 ----------------------------------

func SolvePart2(g *Graph) int {
	startNodeID := getNodeID(0, 1)
	endNodeID := getNodeID(140, 139)
	if runTest {
		endNodeID = getNodeID(22, 21)
	}
	path := []string{}
	maxLength := 0
	dfs(g, g.Nodes[startNodeID], g.Nodes[endNodeID], &map[string]bool{}, &path, &maxLength)
	return maxLength - 1
}

// -------------------------- Main entry -------------------------------------

func main() {
	var input []string

	stopwatch := utils.NewStopwatch()
	stopwatch.Start()

	if runTest {
		input = utils.ReadDataFile(TEST_FILE)
	} else {
		input = utils.ReadDataFile(DATA_FILE)
	}

	graph1 := CreateGraph(input)
	graph2 := CreateGraph_Part2(input)

	result1 := SolvePart1(graph1)
	result2 := SolvePart1(graph2)

	stopwatch.Stop()

	// ---------------------- Print results ----------------------------------
	var elapsedTime time.Duration = stopwatch.GetElapsedTime()
	fmt.Println("Running as test:\t", runTest)
	fmt.Println("Result 1:\t\t", result1)
	fmt.Println("Result 2:\t\t", result2)
	fmt.Println("Elapsed time:\t\t", elapsedTime)
}
