// ---------------------------------------------------------------------------
// Golang solution for Advent of Code 2023 day ..
// https://adventofcode.com/2023/day/..
//
// This was created using copilot to assist me in learning Go.
//
// Scenario:
//
// workflows look like this "ex{x>10:one,m<20:two,a>30:R,A}"
//
//		ex: name of the workflow
//		in {}: the rules
//		x>10:one: if x is greater than 10, then do worflow "one"
//		a>30:R: if the part's a is more than 30, the part is immediately rejected (R)
//		Rule "A": because no other rules matched the part, the part is immediately accepted (A).
//
//	 If part is sent to another workflow it never returns to the original workflow.
//
// machine part categories
//
//	x: Extremely cool looking
//	m: Musical
//	a: Aerodynamic
//	s: Shiny
//
// Part 1: add together all of the ratings for all of the parts that get accepted
// Part 2: find the range of possible inputs for all 4 attributes then multiply out the ranges
// ---------------------------------------------------------------------------
package main

import (
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/cdr74/AdventOfCode2023/utils"
)

// runTest defines whether test.data or actual.data should be used
const runTest bool = false
const TEST_FILE string = "test.data"
const DATA_FILE string = "actual.data"

// -------------------------- Common Data Section ----------------------------

type Condition struct {
	name   string
	op     string
	value  int
	target string
}

type Rule struct {
	name       string
	conditions []Condition
}

type Attribute struct {
	description string
	value       int
}

type Part struct {
	attr []Attribute
}

var parts []Part
var workflow []Rule

func (r Rule) getUnconditionalState() string {
	return r.conditions[len(r.conditions)-1].target
}

func (r Rule) nextState(part Part) string {
	for _, condition := range r.conditions {
		if condition.name == "" {
			// last condition has no attribute
			return condition.target
		}

		for _, att := range part.attr {
			// check first condition against all attributes in part
			if condition.name == att.description {
				switch condition.op {
				case "<":
					if att.value < condition.value {
						return condition.target
					}
				case ">":
					if att.value > condition.value {
						return condition.target
					}
				}
			}
		}
	}
	panic("nextState() not found")
}

func getRuleByName(name string) Rule {
	for _, rule := range workflow {
		if rule.name == name {
			return rule
		}
	}
	panic("getRuleByName() not found")
}

// -------------------------- Common Code Section ----------------------------

// takes ex{x>10:one,m<20:two,a>30:R,A}
// return ex and "x>10:one,m<20:two,a>30:R,A"
func extractRuleName(line string) (string, string) {
	idx := strings.Index(line, "{")
	if idx == 0 {
		panic("Invalid input, no { found}")
	}
	name := line[:idx]        // get name
	line = line[idx+1:]       // remove name
	line = line[:len(line)-1] // remove }
	return name, line
}

func parseInput(input []string) {
	pattern := `([a-zA-Z]+)([<>]+)(\d+):`
	regex := regexp.MustCompile(pattern)

	pattern2 := `([a-zA-Z]+)=(\d+)`
	regex2 := regexp.MustCompile(pattern2)

	var idx int
	var row int
	for row = 0; row < len(input) && input[row] != ""; row++ {
		// parse rules into workflow
		// ex{x>10:one,m<20:two,a>30:R,A}

		line := input[row]
		ruleName, line := extractRuleName(line)
		rule := Rule{name: ruleName} // get name

		sections := strings.Split(line, ",")
		for _, section := range sections {
			// conditions are optional, target is required
			var condition Condition
			if strings.Contains(section, ":") {
				matches := regex.FindStringSubmatch(section)
				condition.name = matches[1]
				condition.op = matches[2]
				condition.value = utils.StringToInt(matches[3])
				idx = strings.Index(section, ":")
				section = section[idx+1:]
			}
			condition.target = section
			rule.conditions = append(rule.conditions, condition)
		}
		workflow = append(workflow, rule)
	}

	for row++; row < len(input); row++ {
		// parse parts
		// {x:1,m:2,a:3,s:4}
		var part Part
		line := input[row][1 : len(input[row])-1] // remove {}
		sections := strings.Split(line, ",")
		for _, section := range sections {
			matches := regex2.FindStringSubmatch(section)
			attribute := Attribute{}
			attribute.description = matches[1]
			attribute.value = utils.StringToInt(matches[2])
			part.attr = append(part.attr, attribute)
		}
		parts = append(parts, part)
	}

}

// -------------------------- Puzzle part 1 ----------------------------------

// returns true if part is accepted in thge workflow
func performWorkflow(part Part) bool {
	var nextState string = "in"
	for nextState != "A" && nextState != "R" {
		rule := getRuleByName(nextState)
		nextState = rule.nextState(part)
	}
	return nextState == "A"
}

func SolvePart1() int {
	var accepted []Part
	for _, part := range parts {
		if performWorkflow(part) {
			accepted = append(accepted, part)
		}
	}

	summ := 0
	for _, part := range accepted {
		for _, att := range part.attr {
			summ += att.value
		}
	}
	return summ
}

// -------------------------- Puzzle part 2 ----------------------------------

type Status int

const (
	X Status = iota
	M
	A
	S
)

type Range struct {
	min int
	max int
}

func (r Range) rangeLength() int {
	return r.max - r.min + 1
}

var statusRanges = map[Status]Range{
	X: {1, 4000},
	M: {1, 4000},
	A: {1, 4000},
	S: {1, 4000},
}

func (s Status) String() string {
	return [...]string{"x", "m", "a", "s"}[s]
}

func cloneStatusRanges(input map[Status]Range) map[Status]Range {
	result := make(map[Status]Range)
	for k, v := range input {
		result[k] = v
	}
	return result
}

func combinationsOfStatusRange(sr map[Status]Range) int {
	mult := 1
	for _, r := range sr {
		mult *= r.rangeLength()
	}
	return mult
}

func statusValid(sr map[Status]Range) bool {
	for _, rng := range sr {
		if rng.min >= rng.max {
			return false
		}
	}
	return true
}

func traverseWorkflow(ruleName string, statusRanges map[Status]Range, result *[]int) {
	fmt.Printf("Processing rule: %v,\tRanges %v,Possibilities\t%d\n", ruleName, statusRanges, combinationsOfStatusRange(statusRanges))

	// accepted
	if ruleName == "A" {
		mult := combinationsOfStatusRange(statusRanges)
		*result = append(*result, mult)
		return
	}

	// rejected, discard
	if ruleName == "R" || !statusValid(statusRanges) {
		fmt.Println("Rejected")
		return
	}

	rule := getRuleByName(ruleName)
	for _, condition := range rule.conditions {
		if condition.name == "" {
			// last condition has no name or attribute
			break
		}

		for status, rng := range statusRanges {
			if condition.name == status.String() {
				tmp := rng
				// adjust range for current condition
				switch condition.op {
				case "<":
					rng.max = condition.value - 1
				case ">":
					rng.min = condition.value + 1
				}
				statusRanges[status] = rng
				traverseWorkflow(condition.target, cloneStatusRanges(statusRanges), result)

				// now inverse this scope change, for next condition to continue
				rng = tmp
				switch condition.op {
				case "<":
					rng.min = condition.value
				case ">":
					rng.max = condition.value
				}
				statusRanges[status] = rng
			}
		}
	}

	// last condition has no attribute; this is an alternative path, we already inversed condition
	traverseWorkflow(rule.getUnconditionalState(), cloneStatusRanges(statusRanges), result)
}

// find the range of possible inputs for all 4 attributes then multiply out the ranges
func SolvePart2() int64 {
	var result []int
	traverseWorkflow("in", statusRanges, &result)
	// add up length of ranges
	var summ int64
	for _, r := range result {
		summ += int64(r)
	}
	return summ
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

	parseInput(input)
	//printRules()
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
