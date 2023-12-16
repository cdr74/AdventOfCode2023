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
// Part 2: credit goes to this guy as i f..up my recursion
//
//	https://pastebin.com/djb8RJ85
//
// ---------------------------------------------------------------------------
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

var memoizedResults map[string]int = make(map[string]int)

func containsValue(list []byte, value byte) bool {
	for i := 0; i < len(list); i++ {
		if list[i] == value {
			return true
		}
	}
	return false
}

func memoizedRecursiveCount(line []byte, sequenceList []int) int {
	//fmt.Printf("memoizedRecursiveCount() - line: %s, sequence: %v\n", string(line), sequenceList)
	memoKey := fmt.Sprintf("%s-%v", string(line), sequenceList)
	if val, ok := memoizedResults[memoKey]; ok {
		return val
	}

	count := recursiveCount(line, sequenceList)
	memoizedResults[memoKey] = count

	return count
}

// consumes line from start and removes processed characters, calls memoizedRecursiveCount to cahce results
func recursiveCount(line []byte, sequenceList []int) int {
	for {
		if len(sequenceList) == 0 {
			if containsValue(line, '#') {
				// invalid combination
				return 0
			} else {
				// no more groups to match, valid combination
				return 1
			}
		}

		// sequenceList is not empty, but line is fully processed
		if len(line) == 0 {
			return 0
		}

		// remove leading '.' and restart loop
		if line[0] == '.' {
			line = []byte(strings.TrimLeft(string(line), "."))
			continue
		}

		if line[0] == '?' {
			// try both options and add them up
			line[0] = '.'
			cnt := memoizedRecursiveCount(line, sequenceList)
			line[0] = '#'
			cnt += memoizedRecursiveCount(line, sequenceList)
			return cnt
		}

		if line[0] == '#' {
			// start of a group, consume it

			// group must not be followed by a '#' and we need enough chars left in line
			if len(line) < sequenceList[0]+1 || line[sequenceList[0]] == '#' {
				return 0
			}

			// check if for length of group if a '.' exists
			for j := 0; j < sequenceList[0]; j++ {
				if line[j] == '.' {
					// premature end of group, invalid
					return 0
				}
			}

			if len(sequenceList) > 1 {
				// skip to after group and after the dot following the group
				line = line[sequenceList[0]+1:]
				sequenceList = sequenceList[1:]
				continue
			}

			line = line[sequenceList[0]:]
			sequenceList = sequenceList[1:]
			continue
		}

		panic("recursiveCount() - bad input line")
	}
}

// -------------------------- Puzzle part 1 ----------------------------------

// find comninations of # sequences for the sequenceList (ints with sequence lengths)
// input "?###???????? 3,2,1" has 10 possible arrangements
// before and after each group of # there must be a '.' or end of string
func SolvePart1(input []string) int {
	var result int = 0

	for _, line := range input {
		sequenceList := getSequenceList(line)
		sequence := line[:strings.Index(line, " ")] + string('.')

		// clear cache
		memoizedResults = make(map[string]int)
		cnt := memoizedRecursiveCount([]byte(sequence), sequenceList)
		fmt.Printf("line: %s, sequence: %v, results: %d\n", sequence, sequenceList, cnt)
		result += cnt
	}
	return result
}

// -------------------------- Puzzle part 2 ----------------------------------

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
	return result + string('.')
}

// multiply all by 5
func SolvePart2(input []string) int {
	var result int = 0

	for _, line := range input {
		sequenceList := getSequenceList(line)
		sequenceList = multiplyList(sequenceList, 5)

		sequence := line[:strings.Index(line, " ")]
		sequence = multiplySequence(sequence, 5)

		// clear cache
		memoizedResults = make(map[string]int)
		cnt := memoizedRecursiveCount([]byte(sequence), sequenceList)
		//fmt.Printf("line: %s, sequence: %v, results: %d\n", sequence, sequenceList, cnt)
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
