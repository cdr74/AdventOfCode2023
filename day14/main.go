// ---------------------------------------------------------------------------
// Golang solution for Advent of Code 2023 day 14
// https://adventofcode.com/2023/day/14
//
// This was created using copilot to assist me in learning Go.
//
// Scenario: input is a field of 2D points. each point can have 3 values
//
//	"O" a ball that can move
//	"#" a wall that cannot be moved through
//	"." an empty space that can be moved through
//
// Part 1:
//
//	   Tilt the field north, let each ball roll as far as possible
//		  without hitting a wall or another ball.
//		  Calculate the "weight" of each ball by multiplying with
//		  filed length (north - south) - distance from north
//
// Part 2:
// ---------------------------------------------------------------------------
package main

import (
	"fmt"
	"hash/fnv"
	"time"

	"github.com/cdr74/AdventOfCode2023/utils"
)

// runTest defines whether test.data or actual.data should be used
const runTest bool = false
const TEST_FILE string = "test.data"
const DATA_FILE string = "actual.data"

// -------------------------- Common Code Section ----------------------------

const WALL = 99
const EMPTY = 0
const BALL = 1

var ROWS int
var COLS int

var cache = make(map[uint32]int)
var field [][]byte

func inputLineToValues(input []string) {
	ROWS = len(input)
	COLS = len(input[0])
	result := make([][]byte, ROWS)
	for i := range result {
		result[i] = make([]byte, COLS)
	}

	for r, line := range input {
		for c, char := range line {
			switch char {
			case 'O':
				result[r][c] = BALL
			case '#':
				result[r][c] = WALL
			case '.':
				result[r][c] = EMPTY
			}
		}
	}
	field = result
}

func HashField() uint32 {
	hash := fnv.New32()
	hash.Write([]byte(fmt.Sprintf("%v", field)))
	return hash.Sum32()
}

func printField() {
	for r := 0; r < ROWS; r++ {
		for c := 0; c < COLS; c++ {
			switch field[r][c] {
			case WALL:
				fmt.Print("#")
			case EMPTY:
				fmt.Print(".")
			case BALL:
				fmt.Print("O")
			}
		}
		fmt.Printf("\n")
	}
	fmt.Printf("\n\n")
}

func tiltNorth() {
	moveCount := 99
	for moveCount > 0 {
		moveCount = 0
		for r := 1; r < ROWS; r++ {
			for c := 0; c < COLS; c++ {
				if field[r][c] == BALL && field[r-1][c] == EMPTY {
					field[r][c] = EMPTY
					field[r-1][c] = BALL
					moveCount++
				}
			}
		}
	}
}

func tiltSouth() {
	moveCount := 99
	for moveCount > 0 {
		moveCount = 0
		for r := ROWS - 2; r >= 0; r-- {
			for c := 0; c < COLS; c++ {
				if field[r][c] == BALL && field[r+1][c] == EMPTY {
					field[r][c] = EMPTY
					field[r+1][c] = BALL
					moveCount++
				}
			}
		}
	}
}

func tiltEast() {
	moveCount := 99
	for moveCount > 0 {
		moveCount = 0
		for c := COLS - 2; c >= 0; c-- {
			for r := 0; r < ROWS; r++ {
				if field[r][c] == BALL && field[r][c+1] == EMPTY {
					field[r][c] = EMPTY
					field[r][c+1] = BALL
					moveCount++
				}
			}
		}
	}
}

func tiltWest() {
	moveCount := 99
	for moveCount > 0 {
		moveCount = 0
		for c := 1; c < COLS; c++ {
			for r := 0; r < ROWS; r++ {
				if field[r][c] == BALL && field[r][c-1] == EMPTY {
					field[r][c] = EMPTY
					field[r][c-1] = BALL
					moveCount++
				}
			}
		}
	}
}

func countFromNorth() int {
	result := 0
	for r := 0; r < ROWS; r++ {
		for c := 0; c < COLS; c++ {
			if field[r][c] == BALL {
				result += (ROWS - r) * 1
			}
		}
	}
	return result
}

// -------------------------- Puzzle part 1 ----------------------------------

func SolvePart1() int {
	tiltNorth()
	result := countFromNorth()
	return result
}

// -------------------------- Puzzle part 2 ----------------------------------

func SolvePart2() int {
	LOOP_COUNT := 1000000000
	doingRest := false
	for i := 0; i < LOOP_COUNT; i++ {
		tiltNorth()
		tiltWest()
		tiltSouth()
		tiltEast()
		hash := HashField()

		if loopStart, ok := cache[hash]; ok && !doingRest {
			loopLength := i - loopStart
			rest := (LOOP_COUNT - loopStart) % loopLength
			i = LOOP_COUNT - loopLength - rest
			doingRest = true
		}
		cache[hash] = i

	}

	result := countFromNorth()
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

	// reset field
	inputLineToValues(input)
	result2 := SolvePart2()

	stopwatch.Stop()

	// ---------------------- Print results ----------------------------------
	var elapsedTime time.Duration = stopwatch.GetElapsedTime()
	fmt.Println("Running as test:\t", runTest)
	fmt.Println("Result 1:\t\t", result1)
	fmt.Println("Result 2:\t\t", result2)
	fmt.Println("Elapsed time:\t\t", elapsedTime)
}
