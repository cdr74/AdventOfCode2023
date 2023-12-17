// ---------------------------------------------------------------------------
// Golang solution for Advent of Code 2023 day 17
// https://adventofcode.com/2023/day/17
//
// This was created using copilot to assist me in learning Go.
//
// Scenario: Get the crucible from the lava pool to the machine parts factory.
//    To do this, you need to minimize heat loss while choosing a route that
//    doesn't require the crucible to go in a straight line for too long.
//
//	Input:
//    lines with number such as "2413432311323"
//    number = single digit that represents the amount of heat loss
//             if the crucible enters that block
//
//  Rule for movement
//     at most three blocks in a single direction
//     then turn left or right
//
// Part 1: Start top left, destination bottom right
// Part 2:
// ---------------------------------------------------------------------------
package main

import (
	"fmt"
	"math"
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

var heat int = math.MaxInt32

// Position represents a move from one cell to another.
type Position struct {
	Row, Col int
}

func isValidMove(row, col int) bool {
	return row >= 0 && row < ROWS && col >= 0 && col < COLS
}

func is4Straight(path []Position, move Position) bool {
	if len(path) < MAX_STRAIGHT {
		return false
	} else {
		current := len(path) - 1
		if move.Col == path[current].Col &&
			path[current].Col == path[current-1].Col &&
			path[current].Col == path[current-2].Col {
			return true
		}
		if move.Row == path[current].Row &&
			path[current].Row == path[current-1].Row &&
			path[current].Row == path[current-2].Row {
			return true
		}
	}
	return false
}

func heatOfPath(path []Position) int {
	var heat int = 0
	for _, move := range path {
		heat += int(field[move.Row][move.Col])
	}
	return heat
}

func getMoves(row, col int, path []Position) []Position {
	var moves []Position
	if len(path) == 1 {
		// no assumption about start position
		moves = append(moves, Position{row + 1, col})
		moves = append(moves, Position{row, col + 1})
		moves = append(moves, Position{row - 1, col})
		moves = append(moves, Position{row, col - 1})
	} else {
		current := len(path) - 1
		if path[current].Col == path[current-1].Col {
			// We were moving vertically, so get horizontal moves
			moves = append(moves, Position{row, col + 1})
			moves = append(moves, Position{row, col - 1})
			if path[current].Row < path[current-1].Row {
				moves = append(moves, Position{row - 1, col})
			} else {
				moves = append(moves, Position{row + 1, col})
			}
		} else if path[current].Row == path[current-1].Row {
			// We were moving horizontally, so get vertical moves
			moves = append(moves, Position{row + 1, col})
			moves = append(moves, Position{row - 1, col})
			if path[current].Col < path[current-1].Col {
				moves = append(moves, Position{row, col - 1})
			} else {
				moves = append(moves, Position{row, col + 1})
			}
		}
	}
	return moves
}

// findAllPaths finds all valid paths without cycles using backtracking.
func findAllPaths(row, col int, visited [][]bool, path []Position, allPaths *[][]Position) {
	currentHeat := heatOfPath(path)
	if currentHeat >= heat {
		// hotter than already existing path, so stop exploring
		return
	}

	if row == ROWS-1 && col == COLS-1 {
		// Reached the destination, add the path to the result and update heat
		if currentHeat <= heat {
			for _, move := range path {
				fmt.Printf("(%d, %d) -> ", move.Row, move.Col)
			}
			fmt.Println(heatOfPath(path))
			heat = currentHeat
			*allPaths = append(*allPaths, append([]Position{}, path...))
		}
		return
	}

	// Mark the current cell as visited
	visited[row][col] = true

	// Get all 3 valid moves, forward, left and right based on where we were last in path
	moves := getMoves(row, col, path)

	for _, move := range moves {
		if isValidMove(move.Row, move.Col) &&
			!visited[move.Row][move.Col] &&
			!is4Straight(path, move) {

			path = append(path, move)
			findAllPaths(move.Row, move.Col, visited, path, allPaths)

			// Backtrack: remove the last move to explore other possibilities
			path = path[:len(path)-1]
		}
	}

	// Mark the current cell as not visited (backtrack)
	visited[row][col] = false
}

// based on field create a graph for use of dijkstra algo
// Graph map[string]map[string]int
// first string is start position (row, col), example "0,0"
// second string map is end positions (row, col), example "0,1"
// int is distance between start and end (in this case the heat loss)
// no more than 3 blocks in a single direction
func buildGraph() {
	graph := make(utils.Graph)
	graph["start"] = make(map[string]int)
	graph["start"]["0-0"] = int(field[0][0])

	visited := make([][]bool, ROWS)
	for i := range visited {
		visited[i] = make([]bool, COLS)
	}

	var allPaths [][]Position
	findAllPaths(0, 0, visited, []Position{{0, 0}}, &allPaths)

	for _, path := range allPaths {
		for _, move := range path {
			fmt.Printf("(%d, %d) -> ", move.Row, move.Col)
		}
		heat := heatOfPath(path)
		fmt.Printf("heat: %d\n", heat)
		fmt.Println("End")
	}
}

func SolvePart1() int {
	buildGraph()
	var result int = 0
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

	inputLineToValues(input)
	printField(0, 0)

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
