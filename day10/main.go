// ----------------------------------------------------------------------------
// https://adventofcode.com/2023/day/10
//
// Analyzes a pipe channel and finds
// Part 1: The point furthest away from the start counted in steps
//         Simplification S is only connected to 2 pipes; no dead ends to analyze
// Part 2:
//
// Start is represented by S and the rest is represented by signs that for mpipe elements
// - = horizontal pipe
// | = vertical pipe
// J = up and left
// L = up and right
// 7 = down and left
// F = down and right
// . = empty space
//
// ----------------------------------------------------------------------------

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

// -------------------------- Common Section ---------------------------------

type Position struct {
	Row int
	Col int
}

var Start Position
var pipeline []Position

func inputToArray(input []string) [][]int {
	width := len(input[0]) + 2
	depth := len(input) + 2
	dataArr := utils.CreateIntArray(width, depth, '.')

	for r, row := range input {
		for c, char := range row {
			dataArr[r+1][c+1] = int(char)
			if char == 'S' {
				Start = Position{Row: r + 1, Col: c + 1}
			}
		}
	}

	return dataArr
}

func printMaze(maze [][]int) {
	fmt.Println("")
	for _, row := range maze {
		for _, char := range row {
			if char == 0 || char == 1 {
				fmt.Printf("%d", char)
			} else {
				switch char {
				case 'J':
					fmt.Print("┘")
				case 'L':
					fmt.Print("└")
				case '7':
					fmt.Print("┐")
				case 'F':
					fmt.Print("┌")
				case 'S':
					fmt.Print("S")
				default:
					fmt.Print(string(char))
				}
			}
		}
		fmt.Println("")
	}
}

// -------------------------- Puzzle part 1 ----------------------------------

// not checking for dead ends!
func nextPosition(maze [][]int, currentPos Position, prevPos Position) Position {
	var nextPos Position
	//fmt.Printf("nextPosition() - current pos: %v, sign: %v\n", currentPos, string(maze[currentPos.Row][currentPos.Col]))

	switch maze[currentPos.Row][currentPos.Col] {
	case '|':
		if prevPos.Row < currentPos.Row {
			nextPos = Position{Row: currentPos.Row + 1, Col: currentPos.Col}
		} else {
			nextPos = Position{Row: currentPos.Row - 1, Col: currentPos.Col}
		}
	case '-':
		if prevPos.Col < currentPos.Col {
			nextPos = Position{Row: currentPos.Row, Col: currentPos.Col + 1}
		} else {
			nextPos = Position{Row: currentPos.Row, Col: currentPos.Col - 1}
		}
	case 'J':
		if prevPos.Col == currentPos.Col {
			nextPos = Position{Row: currentPos.Row, Col: currentPos.Col - 1}
		} else {
			nextPos = Position{Row: currentPos.Row - 1, Col: currentPos.Col}
		}
	case 'L':
		if prevPos.Col == currentPos.Col {
			nextPos = Position{Row: currentPos.Row, Col: currentPos.Col + 1}
		} else {
			nextPos = Position{Row: currentPos.Row - 1, Col: currentPos.Col}
		}
	case '7':
		if prevPos.Col == currentPos.Col {
			nextPos = Position{Row: currentPos.Row, Col: currentPos.Col - 1}
		} else {
			nextPos = Position{Row: currentPos.Row + 1, Col: currentPos.Col}
		}
	case 'F':
		if prevPos.Col == currentPos.Col {
			nextPos = Position{Row: currentPos.Row, Col: currentPos.Col + 1}
		} else {
			nextPos = Position{Row: currentPos.Row + 1, Col: currentPos.Col}
		}
	default:
		panicMsg := fmt.Sprintf("No connection found from: %v\n", currentPos)
		panic(panicMsg)
	}

	return nextPos
}

// Returns count of steps in the pipe connected to Start (S) furthest away from it
func SolvePuzzle1(maze [][]int) int {
	var currentPos Position

	// skipping to find by algo, hardcoded start
	if runTest {
		currentPos = Position{Row: 1, Col: 4}
		//currentPos = Position{Row: 3, Col: 2}
	} else {
		currentPos = Position{Row: 31, Col: 121}
	}

	prevPos := Position{Col: Start.Col, Row: Start.Row}
	nextPos := nextPosition(maze, currentPos, prevPos)
	pipeline = append(pipeline, Start)
	pipeline = append(pipeline, currentPos)
	for nextPos != Start {
		fmt.Printf("Current pos: %v\n", currentPos)
		pipeline = append(pipeline, nextPos)
		prevPos = currentPos
		currentPos = nextPos
		nextPos = nextPosition(maze, currentPos, prevPos)
	}

	return (len(pipeline)) / 2
}

// -------------------------- Puzzle part 2 ----------------------------------

func replacePipelineElements(maze [][]int) {
	for _, position := range pipeline {
		char := maze[position.Row][position.Col]
		switch char {
		case 'J', 'L', '|':
			maze[position.Row][position.Col] = 1
		case '-', '7', 'F', 'S':
			// test and actual input use same direction S
			maze[position.Row][position.Col] = 0
		}
	}
}

func replaceRemainingElements(maze [][]int) {
	for r, row := range maze {
		for c, _ := range row {
			if maze[r][c] != 0 && maze[r][c] != 1 {
				maze[r][c] = '.'
			}
		}
	}
}

func calcCount(maze [][]int, r int, c int) int {
	var count int
	for i := c; i < len(maze[r]); i++ {
		if maze[r][i] == 1 {
			count++
		}
	}
	return count
}

// Calculate fields inside the pipeline, we apply a simple algorithm
//  0. Manually replace S with F in data file
//  1. Replace all pipeline lements as follows (count to right north facing):
//     -, F, 7 = 0
//     |, J, L = 1
//
// 2. Repalce all elements in Maze that are not 0 or 1 with '.'
// 3. Count for each dot whether the summ of 0, 1 elements are even or odd
// 4. If even the "." is outside
// 5. If odd the "." is inside and we increment area counter
// 6. Return area counter
func SolvePuzzle2(maze [][]int) int {
	var result int = 0
	printMaze(maze)
	replacePipelineElements(maze)
	printMaze(maze)
	replaceRemainingElements(maze)
	printMaze(maze)
	for r, row := range maze {
		for c, _ := range row {
			if maze[r][c] == '.' {
				count := calcCount(maze, r, c+1)
				if count%2 == 1 {
					result++
				}
			}
		}
	}
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

	maze := inputToArray(input)

	result1 := SolvePuzzle1(maze)
	result2 := SolvePuzzle2(maze)
	stopwatch.Stop()

	// -------------------------------------
	var elapsedTime time.Duration = stopwatch.GetElapsedTime()
	fmt.Println("Running as test:\t", runTest)
	fmt.Println("Result 1:\t\t", result1)
	fmt.Println("Result 2:\t\t", result2)
	fmt.Println("Elapsed time:\t\t", elapsedTime)
}
