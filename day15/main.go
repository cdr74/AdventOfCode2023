// ---------------------------------------------------------------------------
// Golang solution for Advent of Code 2023 day ..
// https://adventofcode.com/2023/day/..
//
// This was created using copilot to assist me in learning Go.
//
// Scenario:
//
// Part 1:
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
const runTest bool = false
const TEST_FILE string = "test.data"
const DATA_FILE string = "actual.data"

// -------------------------- Common Code Section ----------------------------

// -------------------------- Puzzle part 1 ----------------------------------

func hash(input string) int {
	var result int = 0
	for _, c := range input {
		result += int(c)
		result *= 17
		result %= 256
	}
	return result
}

func SolvePart1(input string) int {
	var result int = 0
	steps := strings.Split(input, ",")
	for _, step := range steps {
		h := hash(step)
		result += h
		fmt.Printf("step: %s, hash: %d\n", step, h)
	}
	return result
}

// -------------------------- Puzzle part 2 ----------------------------------

type Lens struct {
	label string
	focal int
}

type Box struct {
	id     int
	lenses []Lens
}

var boxes = make([]Box, 256)

// anything before '=' or '-' is returned as label
// operator is returned as '=' or '-'
// in case of '=' there is a value that follows
func parseStep(step string) (string, string, int) {
	var label string = ""
	var operator string = ""
	var value int = -1

	for _, c := range step {
		if c == '=' || c == '-' {
			operator = string(c)
			break
		}
		label += string(c)
	}

	if operator == "=" {
		value = utils.StringToInt(step[len(label)+1:])
	}

	return label, operator, value
}

func hasLens(box Box, lens Lens) int {
	for pos, l := range box.lenses {
		if l.label == lens.label {
			return pos
		}
	}
	return -1
}

func SolvePart2(input string) int {
	var result int = 0
	steps := strings.Split(input, ",")
	for _, step := range steps {
		label, operator, value := parseStep(step)
		lens := Lens{label: label, focal: value}
		boxID := hash(label)
		lensPosition := hasLens(boxes[boxID], lens)
		switch operator {
		case "=":
			if lensPosition >= 0 {
				boxes[boxID].lenses[lensPosition].focal = lens.focal
			} else {
				boxes[boxID].lenses = append(boxes[boxID].lenses, lens)
			}
		case "-":
			if lensPosition >= 0 {
				// remove the box with given label from lenses
				boxes[boxID].lenses = append(boxes[boxID].lenses[:lensPosition], boxes[boxID].lenses[lensPosition+1:]...)
			}
			// nothing to do in else case
		}
	}

	lens_focus := make(map[string]int)
	for bID, box := range boxes {
		if len(box.lenses) > 0 {
			for lID, lens := range box.lenses {
				focal := (bID + 1) * (lID + 1) * lens.focal
				lens_focus[lens.label] += focal
			}
		}
	}

	for id, focal := range lens_focus {
		result += focal
		fmt.Printf("focal: %s, %d\n", id, focal)
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

	result1 := SolvePart1(input[0])
	result2 := SolvePart2(input[0])
	stopwatch.Stop()

	// ---------------------- Print results ----------------------------------
	var elapsedTime time.Duration = stopwatch.GetElapsedTime()
	fmt.Println("Running as test:\t", runTest)
	fmt.Println("Result 1:\t\t", result1)
	fmt.Println("Result 2:\t\t", result2)
	fmt.Println("Elapsed time:\t\t", elapsedTime)
}
