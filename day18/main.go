// ---------------------------------------------------------------------------
// Golang solution for Advent of Code 2023 day 18
// https://adventofcode.com/2023/day/18
//
// This was created using copilot to assist me in learning Go.
//
// Scenario:
// The digger starts in a 1 meter cube hole in the ground. They then dig
// the specified number of meters up (U), down (D), left (L), or right (R),
// clearing full 1 meter cubes as they go.
//
// The directions are given as seen from above
// "up" = north
// "right" = east
// "down" = south
// "left" = west
//
// Dig plan:
// Lines such as "R 6 (#70c710)"
// R=right, 6=distance, (#70c710)=color
//
// Part 1: the directions give an outer line, also clear the inner area
//         then calculate the total number of cubes cleared.
// Part 2: color is actually the instructions
//         first five hexadecimal digits = distance as a hexadecimal number
//         last is direction: 0 = R, 1 = D, 2 = L, and 3 = U
// ---------------------------------------------------------------------------
package main

import (
	"fmt"
	"math"
	"strconv"
	"strings"
	"time"

	"github.com/cdr74/AdventOfCode2023/utils"
)

// runTest defines whether test.data or actual.data should be used
const runTest bool = false
const TEST_FILE string = "test.data"
const DATA_FILE string = "actual.data"

// -------------------------- Common Data Section ----------------------------

const EMPTY = 0
const FILLED = 100

type Instruction struct {
	direction string
	distance  int
	color     string
}

type Position struct {
	row int64
	col int64
}

var digPlan []Instruction

var field [][]byte
var ROWS int
var COLS int
var min_row int64 = math.MaxInt64
var min_col int64 = math.MaxInt64
var max_row int64 = 0
var max_col int64 = 0

// -------------------------- Common Code Section ----------------------------

func inputLineToData(input []string) {
	ROWS = 2000
	COLS = 2000
	field = make([][]byte, ROWS)

	for i := range field {
		field[i] = make([]byte, COLS)
	}

	for r, line := range input {
		for c, _ := range line {
			field[r][c] = EMPTY
		}
	}

	for _, line := range input {
		parts := strings.Split(line, " ")
		dir := parts[0]
		dis := utils.StringToInt(parts[1])
		col := parts[2]
		col = strings.Trim(col, " ")
		digPlan = append(digPlan, Instruction{direction: dir, distance: dis, color: col})
	}
}

func printField() {
	for r := min_row; r <= max_row; r++ {
		for c := min_col; c <= max_col; c++ {
			if field[r][c] == EMPTY {
				fmt.Printf(" ")
			} else if field[r][c] == FILLED {
				fmt.Printf("F")
			} else {
				fmt.Printf("%d", field[r][c])
			}
		}
		fmt.Println()
	}
	fmt.Println()
}

// -------------------------- Puzzle part 1 ----------------------------------

func followInstructions() {
	pos := Position{row: 1000, col: 1000}
	field[pos.row][pos.col] = 1

	for _, instruction := range digPlan {
		switch instruction.direction {
		case "U":
			for x := 0; x < instruction.distance; x++ {
				pos.row -= 1
				field[pos.row][pos.col] += 1
			}
			if min_row > pos.row {
				min_row = pos.row
			}
		case "D":
			for x := 0; x < instruction.distance; x++ {
				pos.row += 1
				field[pos.row][pos.col] += 1
			}
			if max_row < pos.row {
				max_row = pos.row
			}
		case "L":
			for x := 0; x < instruction.distance; x++ {
				pos.col -= 1
				field[pos.row][pos.col] += 1
			}
			if min_col > pos.col {
				min_col = pos.col
			}
		case "R":
			for x := 0; x < instruction.distance; x++ {
				pos.col += 1
				field[pos.row][pos.col] += 1
			}
			if max_col < pos.col {
				max_col = pos.col
			}
		}
	}
}

func floodFillStart() Position {
	for r := min_row; r <= max_row; r++ {
		for c := min_col; c <= max_col; c++ {
			if field[r][c] > EMPTY && field[r][c+1] > EMPTY {
				// discard row, hope we have an easy start next row
				break
			}
			if field[r][c] > EMPTY && field[r][c+1] == EMPTY {
				pos := Position{row: r, col: c + 1}
				return pos
			}
		}
	}
	panic("Failed floodFillStart()")
}

func floodFill(pos Position) {
	if field[pos.row][pos.col] != EMPTY {
		return
	}
	field[pos.row][pos.col] = FILLED
	floodFill(Position{row: pos.row - 1, col: pos.col})
	floodFill(Position{row: pos.row + 1, col: pos.col})
	floodFill(Position{row: pos.row, col: pos.col + 1})
	floodFill(Position{row: pos.row, col: pos.col - 1})
}

func countArea() int {
	cnt := 0
	for r := min_row; r <= max_row; r++ {
		for c := min_col; c <= max_col; c++ {
			if field[r][c] > 0 {
				cnt++
			}
		}
	}
	return cnt
}

func SolvePart1() int {
	followInstructions()
	//printField()
	pos := floodFillStart()
	floodFill(pos)
	//printField()
	result := countArea()
	return result
}

// -------------------------- Puzzle part 2 ----------------------------------

var polygon []Position

func updateDigPlanBasedOnColor() {
	for x := 0; x < len(digPlan); x++ {
		instruction := &digPlan[x]
		lStr := instruction.color[2:7]
		dStr := instruction.color[7 : len(instruction.color)-1]
		switch dStr {
		case "0":
			instruction.direction = "R"
		case "1":
			instruction.direction = "D"
		case "2":
			instruction.direction = "L"
		case "3":
			instruction.direction = "U"
		}
		i, _ := strconv.ParseInt(lStr, 16, 64)
		instruction.distance = int(i)
		//fmt.Printf("length %d, dir %s\n", instruction.distance, dStr)
	}
}

func calculatePolygonArea(points []Position) int64 {
	var area int64 = 0
	for x := 0; x < len(points)-1; x++ {
		y := x + 1
		p1 := points[x]
		p2 := points[y]
		area = area + int64(p1.row)*int64(p2.col) - int64(p1.col)*int64(p2.row)
	}
	area = area / 2
	return area
}

// the field will be too big to fit into memory ... let's get smarter
func SolvePart2() int64 {
	updateDigPlanBasedOnColor()
	pos := Position{row: 0, col: 0}
	polygon = append(polygon, pos)
	totalPolygonLength := int64(0)
	for _, instruction := range digPlan {
		switch instruction.direction {
		case "U":
			pos = Position{row: pos.row + int64(instruction.distance), col: pos.col}
		case "D":
			pos = Position{row: pos.row - int64(instruction.distance), col: pos.col}
		case "L":
			pos = Position{row: pos.row, col: pos.col - int64(instruction.distance)}
		case "R":
			pos = Position{row: pos.row, col: pos.col + int64(instruction.distance)}
		}
		totalPolygonLength += int64(instruction.distance)
		polygon = append(polygon, pos)
	}

	result := calculatePolygonArea(polygon)
	result += totalPolygonLength/2 + 1

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

	inputLineToData(input)
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
