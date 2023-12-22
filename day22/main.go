// ---------------------------------------------------------------------------
// Golang solution for Advent of Code 2023 day 22
// https://adventofcode.com/2023/day/22
//
// This was created using copilot to assist me in learning Go.
//
// Scenario:

// Data
// Each line of text represents the position of a single brick at a time
// The position is given as two x,y,z coordinates
// - one for each end of the brick
// - separated by a tilde (~)
// Each brick is made up of a single straight line of cubes
// The whole snapshot is aligned to a three-dimensional cube grid.
//
// A line like 2,2,2~2,2,2 means that both ends of the brick are at the same coordinate
// in other words, that the brick is a single cube.
//
// ground is at z=0
//
// Part 1: let brick fall down until they touch ground or a brick below them
//
//	identify which bricks can be disintegrated (do not support other bricks)
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
const runTest bool = true
const TEST_FILE string = "test.data"
const DATA_FILE string = "actual.data"

// -------------------------- Common Code Section ----------------------------

// -------------------------- Puzzle part 1 ----------------------------------

func SolvePart1() int {
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
