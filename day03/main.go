package main

import (
	"fmt"
	"math"

	"github.com/cdr74/AdventOfCode2023/utils"
)

// the flag runTest defines which data file to read
const runTest bool = false
const TEST_FILE string = "test.data"
const DATA_FILE string = "actual.data"

// ---------------------------------------------------------------------------

// Translates strings like "467..114.#" into an int array.
// Digits remain digits, . = -1; Anything else = -2
// Part 2 Gears are -3
// we add a borer of -1
func inputToArray(input []string) [][]int {
	width := len(input[0]) + 2
	depth := len(input) + 2
	dataArr := utils.CreateIntArray(width, depth, -1)

	for d, row := range input {
		for w, char := range row {
			switch char {
			case '.':
				// already initialized, no action needed
			case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
				dataArr[d+1][w+1] = int(char - '0')
			case '*':
				dataArr[d+1][w+1] = -3
			default:
				dataArr[d+1][w+1] = -2
			}
		}
	}

	return dataArr
}

// ---------------------------------------------------------------------------

// starts with rightmost digit, work to left until no digit found
func getNumber(row []int, w int) (int, int) {
	result := row[w]
	len := 1

	for pos := w - 1; pos >= 1 && row[pos] >= 0; pos-- {
		result += int(math.Pow10(len)) * row[pos]
		len++
	}

	return result, len
}

func isIsolated(data [][]int, d int, w int, len int) bool {
	for y := d - 1; y <= d+1; y++ {
		for x := w - len; x <= w+1; x++ {
			if data[y][x] < -1 {
				return false
			}
		}
	}
	return true
}

func SolvePuzzle1(data [][]int) int {
	var result int = 0

	for d := 1; d < len(data)-1; d++ {
		for w := len(data[0]) - 1; w >= 1; w-- {
			if data[d][w] >= 0 {
				number, len := getNumber(data[d], w)
				if !isIsolated(data, d, w, len) {
					result += number
				}
				w -= len
			}
		}
	}

	return result
}

// ---------------------------------------------------------------------------

// searches start of digits to the right
func findStartOfNumber(row []int, start int) int {
	for idx := start; idx < len(row); idx++ {
		if row[idx] < 0 {
			return idx - 1
		}
	}
	panic("Failed to find start of Number")
}

// Checks 3 fields for a positive number and then identifies the actual number
func nbrsOfRow(row []int, start int) []int {
	var results []int
	a := row[start]
	b := row[start+1]
	c := row[start+2]

	switch {
	case a >= 0 && b >= 0 && c >= 0:
		idx := findStartOfNumber(row, start)
		number, _ := getNumber(row, idx)
		results = append(results, number)
	case a >= 0 && b < 0 && c >= 0:
		idx := start
		number, _ := getNumber(row, idx)
		results = append(results, number)
		idx = findStartOfNumber(row, start+2)
		number2, _ := getNumber(row, idx)
		results = append(results, number2)
	case a >= 0 && b < 0 && c < 0:
		idx := start
		number, _ := getNumber(row, idx)
		results = append(results, number)
	case a < 0 && b < 0 && c >= 0:
		idx := findStartOfNumber(row, start+2)
		number, _ := getNumber(row, idx)
		results = append(results, number)
	case a < 0 && b >= 0 && c >= 0:
		idx := findStartOfNumber(row, start+2)
		number, _ := getNumber(row, idx)
		results = append(results, number)
	case a >= 0 && b >= 0 && c < 0:
		idx := start + 1
		number, _ := getNumber(row, idx)
		results = append(results, number)
	case a < 0 && b >= 0 && c < 0:
		idx := start + 1
		number, _ := getNumber(row, idx)
		results = append(results, number)
	}
	return results
}

func hasTwoNumbers(data [][]int, d int, w int) (bool, int) {
	var numbers []int
	numbers = append(numbers, nbrsOfRow(data[d-1], w-1)...)
	numbers = append(numbers, nbrsOfRow(data[d], w-1)...)
	numbers = append(numbers, nbrsOfRow(data[d+1], w-1)...)
	if len(numbers) == 2 {
		//fmt.Printf("hasTwoNumbers>> %v\n", numbers)
		return true, numbers[0] * numbers[1]
	}
	return false, 0
}

func SolvePuzzle2(data [][]int) int {
	var result int = 0

	for d := 1; d < len(data)-1; d++ {
		for w := 1; w < len(data[0]); w++ {
			if (data[d][w]) == -3 {
				// found a star
				hit, power := hasTwoNumbers(data, d, w)
				if hit {
					//fmt.Printf("Found engine part at %d, %d - %d\n", d, w, power)
					result += power
				}
			}
		}
	}
	return result
}

// ---------------------------------------------------------------------------

func main() {
	var input []string

	stopwatch := utils.NewStopwatch()
	stopwatch.Start()

	if runTest {
		input = utils.ReadDataFile(TEST_FILE)
	} else {
		input = utils.ReadDataFile(DATA_FILE)
	}

	dataArr := inputToArray(input)
	result1 := SolvePuzzle1(dataArr)
	result2 := SolvePuzzle2(dataArr)
	stopwatch.Stop()

	// -------------------------------------
	elapsedTime := stopwatch.GetElapsedTime()
	fmt.Println("Running as test:\t", runTest)
	fmt.Println("Result 1:\t\t", result1)
	fmt.Println("Result 2:\t\t", result2)
	fmt.Println("Elapsed time:\t\t", elapsedTime)
}
