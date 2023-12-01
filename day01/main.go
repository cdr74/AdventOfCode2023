// https://adventofcode.com/2023/day/1
package main

import (
	"fmt"
	"regexp"
	"strconv"
	"time"

	"github.com/cdr74/AdventOfCode2023/utils"
)

// the flag runTest defines which data file to read
const runTest bool = false
const TEST_FILE string = "test.data"
const DATA_FILE string = "actual.data"

// ---------------------------------------------------------------------------

func charToNumber(char rune) int {
	if '0' <= char && char <= '9' {
		num, _ := strconv.Atoi(string(char))
		return num
	}
	return -1
}

func SolvePuzzlePart1(input []string) int {
	summ := 0

	for _, line := range input {
		firstDigit := -1
		lastDigit := -1
		currentDigit := -1
		for _, char := range line {
			currentDigit = charToNumber(char)
			if currentDigit != -1 {
				if firstDigit == -1 {
					firstDigit = currentDigit
				} else {
					lastDigit = currentDigit
				}
			}
		}

		if lastDigit == -1 {
			lastDigit = firstDigit
		}

		lineValue := firstDigit*10 + lastDigit
		// fmt.Printf("Line value: %d\n", lineValue)
		summ += lineValue
	}

	return summ
}

// ---------------------------------------------------------------------------

func findFirstAndLastMatch(input string, patterns []string) (firstMatch string, lastMatch string) {
	posFirstMatch := -1
	posLastMatch := -1
	firstMatch = ""
	lastMatch = ""

	// I am sure there's a smarter way
	for _, pattern := range patterns {
		re := regexp.MustCompile(pattern)
		matches := re.FindAllStringIndex(input, -1)

		// If there are matches, update the first and last match
		if len(matches) > 0 {
			if posFirstMatch == -1 || matches[0][0] < posFirstMatch {
				firstMatch = input[matches[0][0]:matches[0][1]]
				posFirstMatch = matches[0][0]
			}
			if posLastMatch == -1 || matches[len(matches)-1][0] > posLastMatch {
				lastMatch = input[matches[len(matches)-1][0]:matches[len(matches)-1][1]]
				posLastMatch = matches[len(matches)-1][0]
			}
		}
	}

	return firstMatch, lastMatch
}

func getValue(patterns []string, target string) int {
	for i, value := range patterns {
		if value == target {
			if i > 9 {
				return i - 10
			}
			return i
		}
	}
	// Target string not found in the list
	return -1
}

func SolvePuzzlePart2(input []string) int {
	patterns := []string{"zero", "one", "two", "three", "four", "five", "six", "seven", "eight", "nine",
		"0", "1", "2", "3", "4", "5", "6", "7", "8", "9"}
	summ := 0
	for _, line := range input {
		firstMatch, lastMatch := findFirstAndLastMatch(line, patterns)
		firstValue := getValue(patterns, firstMatch)
		lastValue := getValue(patterns, lastMatch)
		fmt.Printf("String: %s - First: %v - Last: %v\n", line, firstValue, lastValue)
		summ += firstValue*10 + lastValue
	}

	return summ
}

func main() {
	var result1 int
	var result2 int

	var input []string

	stopwatch := utils.NewStopwatch()
	stopwatch.Start()

	if runTest {
		input = utils.ReadDataFile(TEST_FILE)
	} else {
		input = utils.ReadDataFile(DATA_FILE)
	}

	result1 = SolvePuzzlePart1(input)
	result2 = SolvePuzzlePart2(input)

	stopwatch.Stop()

	// -------------------------------------
	var elapsedTime time.Duration = stopwatch.GetElapsedTime()
	fmt.Println("Running as test:\t", runTest)
	fmt.Printf("Result:\t\t\t%d\n", result1)
	fmt.Printf("Result:\t\t\t%d\n", result2)
	fmt.Println("Elapsed time:\t\t", elapsedTime)
}
