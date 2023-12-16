// ---------------------------------------------------------------------------
// Golang solution for Advent of Code 2023 day 12
// https://adventofcode.com/2023/day/12
//
// This was created using copilot to assist me in learning Go.
//
// Scenario: read a file with lines like this "?###???????? 3,2,1"
//
//	        . = working spring
//	        # = defect spring
//	        ? = can't tell
//			3,2,1 = continous group of damaged springs
//					(each group is seperated by at least one working spring)
//
// Part 1: count all possible combinations of aarrangements on all lines
// Part 2:
// ---------------------------------------------------------------------------
package main

import (
	"fmt"
	"strings"
	"time"

	"github.com/cdr74/AdventOfCode2023/utils"
)

// runTest defines whether test.data or actual.data should be used
const runTest bool = true
const TEST_FILE string = "test.data"
const DATA_FILE string = "actual.data"

// -------------------------- Common Code Section ----------------------------

func getSequenceList(line string) []int {
	var result []int
	numberList := line[strings.Index(line, " "):]
	numbers := strings.Split(numberList, ",")
	for _, number := range numbers {
		result = append(result, utils.StringToInt(number))
	}
	return result
}

// -------------------------- Puzzle part 1 ----------------------------------

func isValidSequence(line string, idx int, length int) bool {
	if idx+length > len(line) {
		return false
	}

	// check sequence
	for i := idx; i < idx+length; i++ {
		if !(line[i] == '#' || line[i] == '?') {
			return false
		}
	}

	// check right border
	if idx+length < len(line) {
		if !(line[idx+length] == '.' || line[idx+length] == '?') {
			return false
		}
	}

	// check left border
	if idx > 0 {
		if !(line[idx-1] == '.' || line[idx-1] == '?') {
			return false
		}
	}

	return true
}

func createCombinationsHelper(line []rune, index, count, maxCount int, result *[]string) {
	if index == len(line) {
		*result = append(*result, string(line))
		return
	}

	if line[index] == '?' {
		line[index] = '#'
		if count+1 <= maxCount {
			createCombinationsHelper(line, index+1, count+1, maxCount, result)
		}
		line[index] = '.'
		if count <= maxCount {
			createCombinationsHelper(line, index+1, count, maxCount, result)
		}
		line[index] = '?' // backtrack for other combinations
	} else if line[index] == '#' {
		if count+1 <= maxCount {
			createCombinationsHelper(line, index+1, count+1, maxCount, result)
		}
	} else {
		createCombinationsHelper(line, index+1, count, maxCount, result)
	}
}

func createAllCombinations(line string, maxSequenceLength int, maxCount int) []string {
	var result []string
	createCombinationsHelper([]rune(line), 0, 0, maxCount, &result)
	return result
}

func maxOfList(list []int) int {
	var result int = 0
	for _, item := range list {
		if item > result {
			result = item
		}
	}
	return result
}

func countAllCombinations(line string, sequenceList []int) int {
	var result int = 0
	sumOfSequence := utils.SumOfArray(sequenceList)
	lineCombinations := createAllCombinations(line, maxOfList(sequenceList), sumOfSequence)

	seqID := 0
	for _, combination := range lineCombinations {
		for pos := 0; pos < len(combination); pos++ {
			if combination[pos] != '.' {
				if isValidSequence(combination, pos, sequenceList[seqID]) {
					pos += sequenceList[seqID] - 1
					seqID++
					if seqID == len(sequenceList) {
						result++
						break
					}
				} else {
					break
				}
			}
		}
		seqID = 0
	}

	//fmt.Printf("countAllCombinations() - line: %s, sequence: %v, combinations: %d, results: %d\n", line, sequenceList, len(lineCombinations), result)
	return result
}

// find comninations of # sequences for the sequenceList (ints with sequence lengths)
// input "?###???????? 3,2,1" has 10 possible arrangements
// before and after each group of # there must be a '.' or end of string
func SolvePart1(input []string) int {
	var result int = 0

	for _, line := range input {
		sequenceList := getSequenceList(line)
		sequence := line[:strings.Index(line, " ")]
		result += countAllCombinations(sequence, sequenceList)
	}
	return result
}

// -------------------------- Puzzle part 2 ----------------------------------

var sequenceList []int

func multiplyList(list []int, factor int) []int {
	var result []int
	for i := 0; i < factor; i++ {
		result = append(result, list...)
	}
	return result
}

func multiplySequence(sequence string, factor int) string {
	var result string
	for i := 0; i < factor; i++ {
		if i == 0 {
			result += sequence
		} else {
			result += fmt.Sprintf("?%s", sequence)
		}
	}
	return result
}

func getSequenceLength(sequenceList []int, sequenceIdx int) int {
	if sequenceIdx < len(sequenceList) {
		return sequenceList[sequenceIdx]
	}
	return 0
}

/*
func recursiveCombinationCount(line []byte, charPos int, currentSequenceCount int, sequenceIdx int, results *[]string) {
	sequenceLength := getSequenceLength(sequenceList, sequenceIdx)

	for i := charPos; i < len(line); i++ {
		if line[i] == '#' {
			currentSequenceCount++
			if currentSequenceCount > sequenceLength {
				// wrong branch, too many # in sequence; need to backtrack
				return
			}
		} else if line[i] == '.' {
			if currentSequenceCount > 0 && currentSequenceCount < sequenceLength {
				// wrong branch, not enough # in sequence; need to backtrack
				return
			}
			if currentSequenceCount == sequenceLength {
				// found valid sequence, move to next sequence in list
				currentSequenceCount = 0
				sequenceIdx++
				sequenceLength = getSequenceLength(sequenceList, sequenceIdx)
			}
		} else if line[i] == '?' {
			if sequenceLength > 0 && currentSequenceCount < sequenceLength {
				// no point to go this path if we are done with last sequence
				line[i] = '#'
				recursiveCombinationCount(line, i, currentSequenceCount, sequenceIdx, results)
				// undo change for backtracking
				line[i] = '?'
			}
			line[i] = '.'
			recursiveCombinationCount(line, i, currentSequenceCount, sequenceIdx, results)
			// undo change for backtracking
			line[i] = '?'
			// tried both options, backtrack
			return
		}
	}
	// all sequences consumed; we can end with a '.' or a '#'
	if currentSequenceCount == sequenceLength && sequenceIdx >= len(sequenceList)-1 {
		*results = append(*results, string(line))
	}
}
*/

var memoizedResults map[string]int = make(map[string]int)

func memoizedRecursiveCount(line []byte, charPos int, currentSequenceCount int, sequenceIdx int) int {
	memoKey := fmt.Sprintf("%s-%d-%d-%d", string(line), charPos, currentSequenceCount, sequenceIdx)
	fmt.Printf("memoKey: %v\n", memoKey)
	if val, ok := memoizedResults[memoKey]; ok {
		//fmt.Printf("Cache hit: %v\n", memoKey)
		return val
	} else {
		//fmt.Printf("Cache miss: %v\n", memoKey)
	}

	variationCount := 0
	sequenceLength := getSequenceLength(sequenceList, sequenceIdx)

	for i := charPos; i < len(line); i++ {
		if line[i] == '#' {
			currentSequenceCount++
			if currentSequenceCount > sequenceLength {
				// wrong branch, too many # in sequence; need to backtrack
				return variationCount
			}
		} else if line[i] == '.' {
			if currentSequenceCount > 0 && currentSequenceCount < sequenceLength {
				// wrong branch, not enough # in sequence; need to backtrack
				return variationCount
			}
			if currentSequenceCount == sequenceLength {
				// found valid sequence, move to next sequence in list
				currentSequenceCount = 0
				sequenceIdx++
				sequenceLength = getSequenceLength(sequenceList, sequenceIdx)
			}
		} else if line[i] == '?' {
			if sequenceLength > 0 && currentSequenceCount < sequenceLength {
				// no point to go this path if we are done with last sequence
				line[i] = '#'
				variationCount += memoizedRecursiveCount(line, i, currentSequenceCount, sequenceIdx)
				// undo change for backtracking
				line[i] = '?'
			}
			line[i] = '.'
			variationCount += memoizedRecursiveCount(line, i, currentSequenceCount, sequenceIdx)
			// undo change for backtracking
			line[i] = '?'
			// tried both options, backtrack

			memoKey = fmt.Sprintf("%d-%d-%d", i, currentSequenceCount, sequenceIdx)
			memoizedResults[memoKey] = variationCount
			return variationCount
		}
	}
	// all sequences consumed; we can end with a '.' or a '#'
	if currentSequenceCount == sequenceLength && sequenceIdx >= len(sequenceList)-1 {
		variationCount++
		memoKey = fmt.Sprintf("%d-%d-%d", charPos, currentSequenceCount, sequenceIdx)
		memoizedResults[memoKey] = variationCount
	}

	return variationCount
}

// multiply all by 5
func SolvePart2(input []string) int {
	var result int = 0

	for _, line := range input {
		sequenceList = getSequenceList(line)
		// factor 5
		sequenceList = multiplyList(sequenceList, 5)
		sequence := line[:strings.Index(line, " ")]
		sequence = multiplySequence(sequence, 5)
		//results := []string{}
		//recursiveCombinationCount([]byte(sequence), 0, 0, 0, &results)
		//fmt.Printf("line: %s, sequence: %v, results: %d\n", sequence, sequenceList, len(results))
		// result += len(results)

		// clear cache
		memoizedResults = make(map[string]int)
		cnt := memoizedRecursiveCount([]byte(sequence), 0, 0, 0)
		fmt.Printf("line: %s, sequence: %v, results: %d\n", sequence, sequenceList, cnt)
		result += cnt
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

	result1 := SolvePart1(input)
	result2 := SolvePart2(input)
	stopwatch.Stop()

	// ---------------------- Print results ----------------------------------
	var elapsedTime time.Duration = stopwatch.GetElapsedTime()
	fmt.Println("Running as test:\t", runTest)
	fmt.Println("Result 1:\t\t", result1)
	fmt.Println("Result 2:\t\t", result2)
	fmt.Println("Elapsed time:\t\t", elapsedTime)
}
