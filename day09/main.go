package main

import (
	"fmt"
	"strings"
	"time"

	"github.com/cdr74/AdventOfCode2023/utils"
)

// runTest defines whether test.data or actual.data should be used
const runTest bool = false
const TEST_FILE string = "test.data"
const DATA_FILE string = "actual.data"

// -------------------------- Common Section ---------------------------------

func inputLineToValues(input []string) [][]int {
	var result [][]int

	for idx, line := range input {
		elements := strings.Split(line, " ")
		var lineValues []int
		for _, element := range elements {
			value := utils.StringToInt(element)
			lineValues = append(lineValues, value)
		}
		result = append(result, lineValues)
		idx++
	}

	return result
}

func isLineZero(line []int) bool {
	for _, value := range line {
		if value != 0 {
			return false
		}
	}
	return true
}

func lineDiff(line []int) []int {
	result := make([]int, len(line)-1)

	for idx := 1; idx < len(line); idx++ {
		result[idx-1] = line[idx] - line[idx-1]
	}

	return result
}

// returns an [][]int with the diffs between each value in the line
// down to the line that has only 0 values
func calculateRows(line []int) [][]int {
	var temp [][]int
	temp = append(temp, line)

	for !isLineZero(temp[len(temp)-1]) {
		temp = append(temp, lineDiff(temp[len(temp)-1]))
	}

	return temp
}

func sumValues(results []int) int {
	result := 0
	for _, value := range results {
		result += value
	}
	return result
}

// -------------------------- Puzzle part 1 ----------------------------------

// add last value of current row and last value of row below and append to the current row
func addLastValues(temp [][]int) {
	for idx := len(temp) - 2; idx >= 0; idx-- {
		value := temp[idx+1][len(temp[idx+1])-1] + temp[idx][len(temp[idx])-1]
		temp[idx] = append(temp[idx], value)
	}
}

// for each array of int, create a new int array with the difference between each value
// continue adding int arrays until all values are 0
// then add up the last value of the lowest array to the last value of the array
// above and add that value to the array above
// return the summ of all added last values
func SolvePuzzle1(input [][]int) int {
	var results []int

	for _, line := range input {
		temp := calculateRows(line)
		addLastValues(temp)
		results = append(results, temp[0][len(temp[0])-1])
	}

	return sumValues(results)
}

// -------------------------- Puzzle part 2 ----------------------------------

// add first value of current row and subtract first value of row below and add this first to the current row
func addFirstValues(temp [][]int) {
	for idx := len(temp) - 2; idx >= 0; idx-- {
		value := temp[idx][0] - temp[idx+1][0]
		temp[idx] = append([]int{value}, temp[idx]...)
	}
}

// for each array of int, create a new int array with the difference between each value
// continue adding int arrays until all values are 0
// then add to the left the value that equals the leftmost value below plus the leftmost value on current line
// return the summ of all added last values (row 0)
func SolvePuzzle2(input [][]int) int {
	var results []int

	for _, line := range input {
		temp := calculateRows(line)
		addFirstValues(temp)
		results = append(results, temp[0][0])
	}

	return sumValues(results)
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
	result2 := SolvePuzzle2(values)
	stopwatch.Stop()

	// -------------------------------------
	var elapsedTime time.Duration = stopwatch.GetElapsedTime()
	fmt.Println("Running as test:\t", runTest)
	fmt.Println("Result 1:\t\t", result1)
	fmt.Println("Result 2:\t\t", result2)
	fmt.Println("Elapsed time:\t\t", elapsedTime)
}
