// ---------------------------------------------------------------------------
// Golang solution for Advent of Code 2023 day 21
// https://adventofcode.com/2023/day/21
//
// This was created using copilot to assist me in learning Go.
//
// Scenario:
//
//	We're in a garden at position S, and we want to walk a given number
//	of steps. We can tiles in the garden marked as '.', but can not step on
//	tiles marked as '#'. What we visited we mark as 'O'.
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

func getStartPos(g [][]byte) Position {
	var pos Position
	for i := 0; i < ROWS; i++ {
		for j := 0; j < COLS; j++ {
			if g[i][j] == 'S' {
				pos.row = i
				pos.col = j
				return pos
			}
		}
	}
	panic("Start not found")
}

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

func nextLevelPositions(g [][]byte, queue []Position) []Position {
	nextLevelMap := make(map[Position]struct{})

	for x := 0; x < len(queue); x++ {
		currPos := queue[x]

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
					nextLevelMap[neighbor] = struct{}{}
				}
			}
		} else {
			panic("Illegal state")
		}

	}

	var nextLevel []Position
	for pos := range nextLevelMap {
		nextLevel = append(nextLevel, pos)
	}
	return nextLevel
}

func bfsGardenWalk(g [][]byte, pos Position, depth int) [][]byte {
	var queue []Position

	queue = append(queue, pos)
	for len(queue) > 0 && depth > 0 {
		queue = nextLevelPositions(g, queue)
		depth--
	}

	// at this level our queue holds the actual positions, mark it
	for len(queue) > 0 {
		currPos := queue[0]
		queue = queue[1:]
		g[currPos.row][currPos.col] = 'O'
	}

	return g
}

// -------------------------- Puzzle part 1 ----------------------------------

var stepsPart1 int = 64 // 6 for test data, 64 for actual data

func SolvePart1() int {
	pos := getStartPos(garden)
	g := bfsGardenWalk(garden, pos, stepsPart1)
	//printGarden(g)
	result := countVisitedTiles(g)
	return result
}

// -------------------------- Puzzle part 2 ----------------------------------

// altough passed by reference go copies the a slice if it gets expanded
// only works on uneven multiply for start point calc
func expandMap(g [][]byte, factor int) ([][]byte, Position) {
	pos := getStartPos(g)
	g[pos.row][pos.col] = '.'

	// multiply length of each row by factor
	for r := 0; r < ROWS; r++ {
		row := g[r]
		for f := 0; f < factor-1; f++ {
			g[r] = append(g[r], row...)
		}
	}
	// append rows by factor
	for j := 0; j < factor-1; j++ {
		for y := 0; y < ROWS; y++ {
			g = append(g, g[y])
		}
	}

	// update start position
	half := factor / 2
	pos.row = pos.row + half*ROWS
	pos.col = pos.col + half*COLS
	g[pos.row][pos.col] = 'S'

	ROWS = len(g)
	COLS = len(g[0])

	return g, pos
}

func copy2DSlice(src [][]byte) [][]byte {
	dst := make([][]byte, len(src))
	for i, row := range src {
		newRow := make([]byte, len(row))
		copy(newRow, row)
		dst[i] = newRow
	}
	return dst
}

func SolvePart2() int {
	full := ROWS
	half := full / 2

	largeMap, pos := expandMap(garden, 5)
	largeMap_clean := copy2DSlice(largeMap)

	g1 := bfsGardenWalk(copy2DSlice(largeMap_clean), pos, half)
	t1 := countVisitedTiles(g1)
	fmt.Printf("t1: %d\n", t1)

	g2 := bfsGardenWalk(copy2DSlice(largeMap_clean), pos, half+full)
	t2 := countVisitedTiles(g2)
	fmt.Printf("t2: %d\n", t2)

	g3 := bfsGardenWalk(copy2DSlice(largeMap_clean), pos, half+2*full)
	t3 := countVisitedTiles(g3)
	fmt.Printf("t3: %d\n", t3)

	// with help from reddit - extrapolate with
	// Lagrange's Interpolation formula
	a := (t3 + t1 - 2*t2) / 2
	b := t2 - t1 - a
	c := t1
	n := 26501365 / full
	result := a*n*n + b*n + c

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

	inputToGarden(input)
	result2 := SolvePart2()
	stopwatch.Stop()

	// ---------------------- Print results ----------------------------------
	var elapsedTime time.Duration = stopwatch.GetElapsedTime()
	fmt.Println("Running as test:\t", runTest)
	fmt.Println("Result 1:\t\t", result1)
	fmt.Println("Result 2:\t\t", result2)
	fmt.Println("Elapsed time:\t\t", elapsedTime)
}
