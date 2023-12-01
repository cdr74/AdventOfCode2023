// https://adventofcode.com/2023/day/1
package main

import (
	"fmt"
	"regexp"
	"strconv"
	"time"
	"unicode"

	"github.com/cdr74/AdventOfCode2023/utils"
)

// the flag runTest defines which data file to read
const runTest bool = true
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

// Finds the first occurance of a pattern in the input and returns the location and the index.
// In case none of the pattern matches it returns -1
func findPatternIndexInString(input string, patterns []string) (location int, patternIndex int) {
	patternIndex = -1
	for i, pattern := range patterns {
		re := regexp.MustCompile(pattern)
		matchIndex := re.FindStringIndex(input)
		if matchIndex != nil && (patternIndex == -1 || matchIndex[0] < location) {
			location = matchIndex[0]
			patternIndex = i
		}
	}
	return location, patternIndex
}

// Takes a string such as "zoneight" containing numbers as text that might overlap and other random
// characters and returns a list of the numbers found in order [1 8]
func splitStringByTextNumbers(input string) []byte {
	var numbers = []string{"zero", "one", "two", "three", "four", "five", "six", "seven", "eight", "nine"}
	var result []byte

	for i := 0; i < len(input); i++ {
		subString := input[i:]
		location, matchIndex := findPatternIndexInString(subString, numbers)
		if matchIndex == -1 {
			break
		} else {
			result = append(result, byte(matchIndex))
			i += location
		}
	}
	return result
}

// Takes a string such as "xtwone3four" containing numbers as text that might overlap or digits and other random
// characters and returns a list of the numbers found in order [2 1 3 4]
func splitStringWithDigitsAndText(input string) []byte {
	var result []byte
	current := ""

	for _, char := range input {
		if unicode.IsDigit(char) {
			if current != "" {
				// split the string before the digit further into numbers[] if required
				numbersOrChars := splitStringByTextNumbers(current)
				result = append(result, numbersOrChars...)
				current = ""
			}
			// add the actual digit
			result = append(result, byte(charToNumber(char)))
		} else {
			current += string(char)
		}
	}

	if current != "" {
		numbersOrChars := splitStringByTextNumbers(current)
		result = append(result, numbersOrChars...)
	}

	return result
}

func SolvePuzzlePart2(input []string) int {
	summ := 0
	for _, line := range input {
		digitList := splitStringWithDigitsAndText(line)

		var firstDigit byte = 0
		var lastDigit byte = 0
		if len(digitList) >= 1 {
			firstDigit = digitList[0]
		}
		if len(digitList) >= 2 {
			lastDigit = digitList[len(digitList)-1]
		}
		if lastDigit == 0 {
			lastDigit = firstDigit
		}

		lineValue := firstDigit*10 + lastDigit
		fmt.Printf("Line: %s - digitList: %v - Value: %v\n", line, digitList, lineValue)
		summ += int(lineValue)
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
