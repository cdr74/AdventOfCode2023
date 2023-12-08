package main

import (
	"fmt"
	"strings"
	"time"

	"github.com/cdr74/AdventOfCode2023/utils"
)

// the flag runTest defines which data file to read test.data or actual.data
const runTest bool = false
const TEST_FILE string = "test.data"
const DATA_FILE string = "actual.data"

// ---------------------------------------------------------------------------

type Transition struct {
	Position string
	Left     string
	Right    string
}

// Input format: AAA = (BBB, BBB)
func parseTransition(input string) Transition {
	var transition Transition

	// Extract the position
	positionEnd := strings.Index(input, " =")
	transition.Position = input[:positionEnd]

	// Extract the left and right elements
	leftStart := strings.Index(input, "(") + 1
	rightEnd := strings.Index(input, ")")
	elements := strings.Split(input[leftStart:rightEnd], ", ")
	transition.Left = elements[0]
	transition.Right = elements[1]

	return transition
}

// ---------------------------------------------------------------------------

// instructions are a sequence of L, R indicating whether a transition leads to left or right next
// transitions are positions with a next point to reach based on an instruction
// instructions are followed till a transition leads to ZZZ, if end of instructions is reached
// start them from the beginning
func SolvePuzzle1(instructions string, transitions map[string]Transition) int {
	var result int = 0
	instructionList := []rune(instructions)
	idx := 0

	currentPosition := "AAA"
	for currentPosition != "ZZZ" {
		transition := transitions[currentPosition]
		if instructionList[idx] == 'L' {
			currentPosition = transition.Left
		} else {
			currentPosition = transition.Right
		}
		idx++
		if idx == len(instructionList) {
			idx = 0
		}
		result++
	}

	return result
}

// ---------------------------------------------------------------------------

// gcd calculates the greatest common divisor using Euclid's algorithm
func gcd(a, b uint64) uint64 {
	for b != 0 {
		a, b = b, a%b
	}
	return a
}

// lcm calculates the least common multiple of two integers
func lcm(a, b uint64) uint64 {
	return a / gcd(a, b) * b
}

// lcmArray calculates the least common multiple of an array of integers
func lcmArray(numbers []uint64) uint64 {
	result := numbers[0]
	for i := 1; i < len(numbers); i++ {
		result = lcm(result, numbers[i])
	}
	return result
}

// more complex version of SolvePuzzle1
// - we have multiple start points, namely each Position that ends in A, eg ZBA
// - we follow the instructions simultanously for each start point
// - we stop when all start points reach a position that ends in Z, eg BQZ
// instructions are a sequence of L, R indicating whether a transition leads to left or right next
// transitions are positions with a next point to reach based on an instruction
func SolvePuzzle2(instructions string, transitions map[string]Transition) uint64 {
	var result uint64 = 1
	instructionList := []rune(instructions)

	// initialize the start positions
	var positionList []string
	for position := range transitions {
		if position[len(position)-1] == 'A' {
			positionList = append(positionList, position)
		}
	}

	lengthOfPath := make([]uint64, len(positionList))
	for x, position := range positionList {
		idx := 0
		length := uint64(0)
		currentPosition := position
		for currentPosition[len(currentPosition)-1] != 'Z' {
			transition := transitions[currentPosition]
			if instructionList[idx] == 'L' {
				currentPosition = transition.Left
			} else {
				currentPosition = transition.Right
			}
			idx++
			if idx == len(instructionList) {
				idx = 0
			}
			length++
		}
		lengthOfPath[x] = length
	}

	result = lcmArray(lengthOfPath)

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

	instructions := input[0]
	input = input[2:]

	transitions := make(map[string]Transition)
	for line := range input {
		transition := parseTransition(input[line])
		transitions[transition.Position] = transition
	}

	// new test data invalidates puzzle 1
	//result1 := SolvePuzzle1(instructions, transitions)
	result2 := SolvePuzzle2(instructions, transitions)
	stopwatch.Stop()

	// -------------------------------------
	var elapsedTime time.Duration = stopwatch.GetElapsedTime()
	fmt.Println("Running as test:\t", runTest)
	//fmt.Println("Result 1:\t\t\t", result1)
	fmt.Println("Result 2:\t\t\t", result2)
	fmt.Println("Elapsed time:\t\t", elapsedTime)
}
