// ---------------------------------------------------------------------------
// Golang solution for Advent of Code 2023 day 21
// https://adventofcode.com/2023/day/21
//
// This was created using copilot to assist me in learning Go.
//
// Scenario:
//   We're in a garden at position S, and we want to walk a given number
//   of steps. We can tiles in the garden marked as '.', but can not step on
//   tiles marked as '#'. What we visited we mark as 'O'.
//
// Part 1: How many tiles did we visit with 64 steps
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

var garden [][]byte
var ROWS int
var COLS int

type Position struct {
	row int
	col int
}

func inputToGarden(input []string) {
	ROWS = len(input)
	COLS = len(input[0])
	garden = make([][]byte, ROWS)
	for i := 0; i < ROWS; i++ {
		garden[i] = make([]byte, COLS)
		for j := 0; j < COLS; j++ {
			garden[i][j] = input[i][j]
		}
	}
}

func printGarden(g [][]byte) {
	for i := 0; i < ROWS; i++ {
		for j := 0; j < COLS; j++ {
			fmt.Printf("%c", g[i][j])
		}
		fmt.Printf("\n")
	}
	fmt.Printf("\n\n")
}

func getStartPos() Position {
	var pos Position
	for i := 0; i < ROWS; i++ {
		for j := 0; j < COLS; j++ {
			if garden[i][j] == 'S' {
				pos.row = i
				pos.col = j
				return pos
			}
		}
	}
	return pos
}

// -------------------------- Puzzle part 1 ----------------------------------

func countVisitedTiles(g [][]byte) int {
	var result int = 0
	for i := 0; i < ROWS; i++ {
		for j := 0; j < COLS; j++ {
			if g[i][j] == 'O' {
				result++
			}
		}
	}
	return result
}

// creates a copy of the passed 2d array
func copyGarden(g [][]byte) [][]byte {
	copy := make([][]byte, len(g))
	for i := 0; i < len(g); i++ {
		copy[i] = make([]byte, len(g[i]))
		for j := 0; j < len(copy[i]); j++ {
			copy[i][j] = g[i][j]
		}
	}
	return copy
}

func containsPosition(slice []Position, target Position) bool {
	for _, pos := range slice {
		if pos == target {
			return true
		}
	}
	return false
}

func isValidPosition(pos Position) bool {
	return pos.row >= 0 && pos.row < ROWS && pos.col >= 0 && pos.col < COLS
}

func walk(g [][]byte, queue []Position) []Position {
	var nextLevel []Position

	for len(queue) > 0 {
		currPos := queue[0]
		queue = queue[1:]

		// Check if the current position is a valid tile to visit
		if g[currPos.row][currPos.col] != '#' {

			// Add neighboring positions to the queue
			neighbors := []Position{
				{currPos.row - 1, currPos.col}, // Up
				{currPos.row + 1, currPos.col}, // Down
				{currPos.row, currPos.col - 1}, // Left
				{currPos.row, currPos.col + 1}, // Right
			}

			for _, neighbor := range neighbors {
				if isValidPosition(neighbor) && g[neighbor.row][neighbor.col] != '#' {
					// check if nextLevel already contains neighbor
					if !containsPosition(nextLevel, neighbor) {
						nextLevel = append(nextLevel, neighbor)
					}
				}
			}
		} else {
			panic("Illegal state")
		}

	}
	return nextLevel
}

func bfsGardenWalk(g [][]byte, pos Position, depth int) [][]byte {
	var queue []Position

	queue = append(queue, pos)
	for len(queue) > 0 && depth > 0 {
		nextLevelQueue := walk(g, queue)
		fmt.Printf("nextLevelQueue %v\n", nextLevelQueue)
		queue = nextLevelQueue
		depth--
	}

	// mark last level
	for len(queue) > 0 {
		currPos := queue[0]
		queue = queue[1:]
		g[currPos.row][currPos.col] = 'O' // Mark as visited
	}

	return g
}

var stepsPart1 int = 64 // 6 for test data, 64 for actual data

func SolvePart1() int {
	pos := getStartPos()
	fmt.Println("Start position:", pos)
	printGarden(garden)
	g := bfsGardenWalk(copyGarden(garden), pos, stepsPart1)
	printGarden(g)
	result := countVisitedTiles(g)
	return result
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

	inputToGarden(input)
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
