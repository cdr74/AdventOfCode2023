// ---------------------------------------------------------------------------
// https://adventofcode.com/2023/day/11
//
// Reading a star map from text file where "." are space and "#" represent stars
// empty columns and ampty rows need to be doubled (space exmasion)
//
// Puzzle 1: calculate the distance between each pair of stars n(n-1)/2
// Puzzle 2: adds a twist to space expansion; every empty col / row is expanded
//           by 1000000 empty cols/rows making a representation as matrix imposssible
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

// -------------------------- Common Section ---------------------------------

type Position struct {
	row int
	col int
}

func inputLineToValues(input []string) [][]int {
	var result [][]int

	for _, line := range input {
		var row []int
		for _, char := range line {
			if char == '.' {
				row = append(row, 0)
			} else {
				row = append(row, 1)
			}
		}
		result = append(result, row)
	}
	return result
}

func isColumnEmpty(starMap [][]int, idx int) bool {
	for _, row := range starMap {
		if row[idx] == 1 {
			return false
		}
	}
	return true
}

func isRowEmpty(row []int) bool {
	for _, value := range row {
		if value == 1 {
			return false
		}
	}
	return true
}

func getStarList(starMap [][]int) []Position {
	var result []Position
	for r := 0; r < len(starMap); r++ {
		for c := 0; c < len(starMap[r]); c++ {
			if starMap[r][c] == 1 {
				result = append(result, Position{row: r, col: c})
			}
		}
	}
	return result
}

func getListofEmptyRows(starMap [][]int) []int {
	var result []int
	for r := 0; r < len(starMap); r++ {
		if isRowEmpty(starMap[r]) {
			result = append(result, r)
		}
	}
	return result
}

func getListofEmptyColumns(starMap [][]int) []int {
	var result []int
	for c := 0; c < len(starMap[0]); c++ {
		if isColumnEmpty(starMap, c) {
			result = append(result, c)
		}
	}
	return result
}

// returns how many elements are in a list between start and end
func countOfValuesBetween(start int, end int, values []int) int {
	var result int = 0
	for _, value := range values {
		if value > start && value < end {
			result++
		}
	}
	return result
}

func getDistanceBetweenStars(star1 Position, star2 Position, emptyRows []int, emptyCols []int, expansionFactor int) int {
	xFrom, xTo := utils.Min(star1.col, star2.col), utils.Max(star1.col, star2.col)
	yFrom, yTo := utils.Min(star1.row, star2.row), utils.Max(star1.row, star2.row)

	xDist := utils.Abs(xTo-xFrom) - countOfValuesBetween(xFrom, xTo, emptyCols) + countOfValuesBetween(xFrom, xTo, emptyCols)*expansionFactor
	yDist := utils.Abs(yTo-yFrom) - countOfValuesBetween(yFrom, yTo, emptyRows) + countOfValuesBetween(yFrom, yTo, emptyRows)*expansionFactor

	return xDist + yDist
}

func getDistanceBetweenAllStars(starMap [][]int, expansionFactor int) int {
	var result int = 0

	emptyRows := getListofEmptyRows(starMap)
	emptyCols := getListofEmptyColumns(starMap)
	starList := getStarList(starMap)

	for i := 0; i < len(starList); i++ {
		for j := i + 1; j < len(starList); j++ {
			result += getDistanceBetweenStars(starList[i], starList[j], emptyRows, emptyCols, expansionFactor)
		}
	}
	return result
}

// -------------------------- Puzzle part 1 ----------------------------------

func SolvePuzzle1(starMap [][]int) int {
	return getDistanceBetweenAllStars(starMap, 2)
}

// -------------------------- Puzzle part 2 ----------------------------------

func SolvePuzzle2(starMap [][]int) int {
	return getDistanceBetweenAllStars(starMap, 1000000)
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

	values := inputLineToValues(input)
	result1 := SolvePuzzle1(values)

	values = inputLineToValues(input)
	result2 := SolvePuzzle2(values)

	stopwatch.Stop()

	// -------------------------------------
	var elapsedTime time.Duration = stopwatch.GetElapsedTime()
	fmt.Println("Running as test:\t", runTest)
	fmt.Println("Result 1:\t\t", result1)
	fmt.Println("Result 2:\t\t", result2)
	fmt.Println("Elapsed time:\t\t", elapsedTime)
}
