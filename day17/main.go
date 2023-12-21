// ---------------------------------------------------------------------------
// Golang solution for Advent of Code 2023 day 17
// https://adventofcode.com/2023/day/17
//
// This was created using copilot to assist me in learning Go.
//
// Scenario: Get the crucible from the lava pool to the machine parts factory.
//
//	   To do this, you need to minimize heat loss while choosing a route that
//	   doesn't require the crucible to go in a straight line for too long.
//
//		Input:
//	   lines with number such as "2413432311323"
//	   number = single digit that represents the amount of heat loss
//	            if the crucible enters that block
//
//	 Rule for movement
//	    at most three blocks in a single direction
//	    then turn left or right
//
// Part 1: Start top left, destination bottom right
// Part 2:
// ---------------------------------------------------------------------------
package main

import (
	"fmt"
	"math"
	"strings"
	"time"

	"github.com/cdr74/AdventOfCode2023/utils"
)

// runTest defines whether test.data or actual.data should be used
const runTest bool = true
const TEST_FILE string = "test.data"
const DATA_FILE string = "actual.data"

// -------------------------- Common Data Section ----------------------------

const MAX_STRAIGHT = 3

var field [][]byte
var ROWS int
var COLS int

// -------------------------- Common Code Section ----------------------------

func inputLineToValues(input []string) {
	ROWS = len(input)
	COLS = len(input[0])
	field = make([][]byte, ROWS)
	for i := range field {
		field[i] = make([]byte, COLS)
	}

	for r, line := range input {
		for c, char := range line {
			field[r][c] = byte(char) - byte('0')
		}
	}
}

func printField(r_crucible int, c_crucible int) {
	for r := 0; r < ROWS; r++ {
		for c := 0; c < COLS; c++ {
			if r == r_crucible && c == c_crucible {
				fmt.Print(" ")
			} else {
				fmt.Print(field[r][c])
			}
		}
		fmt.Println()
	}
	fmt.Println()
}

// -------------------------- Puzzle part 1 ----------------------------------

func printFieldWithDistances(dist map[string]int) {
	for r := 0; r < ROWS; r++ {
		for c := 0; c < COLS; c++ {
			key := fmt.Sprintf("%d-%d", r, c)
			if dist[key] == math.MaxInt32 {
				fmt.Printf("XXXXX\t")
			} else {
				fmt.Printf("%d\t", dist[key])
			}
		}
		fmt.Println()
	}
	fmt.Println()
}

type node struct {
	key  string
	cost int
}

type Graph map[string]map[string]int

func (g Graph) Path(start, target string) (path []string, cost int, err error) {
	explored := make(map[string]bool)   // set of nodes we already explored
	frontier := utils.NewQueue()        // queue of the nodes to explore
	previous := make(map[string]string) // previously visited node

	frontier.Set(start, 0)

	// run until we visited every node in the frontier
	for !frontier.IsEmpty() {
		aKey, aPriority := frontier.Next()
		n := node{aKey, aPriority}

		if n.key == target {
			cost = n.cost
			nKey := n.key
			for nKey != start {
				path = append(path, nKey)
				nKey = previous[nKey]
			}
			break
		}
		explored[n.key] = true

		// loop all the neighboring nodes
		for nKey, nCost := range g[n.key] {
			if explored[nKey] {
				continue
			}

			if _, ok := frontier.Get(nKey); !ok {
				previous[nKey] = n.key
				frontier.Set(nKey, n.cost+nCost)
				continue
			}

			// if violating rule, skip this node as candidate
			if isViolatingRule(nKey, previous) {
				fmt.Printf("Skipping node: %v\n", nKey)
				delete(previous, nKey)
				explored[nKey] = true
				continue
			}

			frontierCost, _ := frontier.Get(nKey)
			nodeCost := n.cost + nCost

			if nodeCost < frontierCost {
				previous[nKey] = n.key
				frontier.Set(nKey, nodeCost)
			}
		}

	}

	// add the origin at the end of the path
	path = append(path, start)

	// reverse the path because it was popilated
	// in reverse, form target to start
	for i, j := 0, len(path)-1; i < j; i, j = i+1, j-1 {
		path[i], path[j] = path[j], path[i]
	}

	return
}

// no n moves in straight line
func isViolatingRule(node string, prev map[string]string) bool {
	var keys []string
	first_col := -1
	first_row := -1
	cnt1 := 0
	cnt2 := 0

	currentKey := node
	for i := 0; i <= MAX_STRAIGHT; i++ {
		keys = append(keys, currentKey)
		if len(currentKey) < 3 {
			return false
		}

		idx := strings.Index(currentKey, "-")
		col := utils.StringToInt(currentKey[:idx])
		row := utils.StringToInt(currentKey[idx+1:])

		if i == 0 {
			first_col = col
			first_row = row
		}
		if first_row == row {
			cnt1++
		}
		if first_col == col {
			cnt2++
		}

		currentKey = prev[currentKey]
	}

	if cnt1 > MAX_STRAIGHT || cnt2 > MAX_STRAIGHT {
		fmt.Printf("Found violation at node: %v, %v, %d, %d\n", node, keys, cnt2, cnt1)
		return true
	}

	return false
}

// type Graph map[string]map[string]int
//
//	Graph{
//		"0-0": {"0-1": 10, "1-0": 20},
//		"0-1": {"0-2": 50},
//	}
func buildGraph() Graph {
	graph := make(map[string]map[string]int)

	for r := 0; r < ROWS; r++ {
		for c := 0; c < COLS; c++ {
			key := fmt.Sprintf("%d-%d", r, c)
			graph[key] = make(map[string]int)

			if r > 0 {
				graph[key][fmt.Sprintf("%d-%d", r-1, c)] = int(field[r-1][c])
			}
			if c > 0 {
				graph[key][fmt.Sprintf("%d-%d", r, c-1)] = int(field[r][c-1])
			}
			if r < ROWS-1 {
				graph[key][fmt.Sprintf("%d-%d", r+1, c)] = int(field[r+1][c])
			}
			if c < COLS-1 {
				graph[key][fmt.Sprintf("%d-%d", r, c+1)] = int(field[r][c+1])
			}
		}
	}

	return graph
}

func SolvePart1() int {
	graph := buildGraph()

	// unmodified dijkstra
	path, cost, _ := graph.Path("0-0", "12-12")
	fmt.Printf("Path from 'start' to '12-12' with lowest cost: %v cost: %v\n", path, cost)

	return cost
}

// -------------------------- Puzzle part 2 ----------------------------------

func SolvePart2() int {
	var result int = 0

	return result
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

	inputLineToValues(input)

	result1 := SolvePart1()
	result2 := SolvePart2()
	stopwatch.Stop()

	// ---------------------- Print results ----------------------------------
	var elapsedTime time.Duration = stopwatch.GetElapsedTime()
	fmt.Println("Running as test:\t", runTest)
	fmt.Println("Result 1:\t\t", result1)
	fmt.Println("Result 2:\t\t", result2)
	fmt.Println("Elapsed time:\t\t", elapsedTime)
}
