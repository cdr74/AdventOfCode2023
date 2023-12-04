package main

import (
	"fmt"
	"strings"
	"time"

	"github.com/cdr74/AdventOfCode2023/utils"
)

// the flag runTest defines which data file to read
const runTest bool = false
const TEST_FILE string = "test.data"
const DATA_FILE string = "actual.data"

// ---------------------------------------------------------------------------

type Ticket struct {
	ID           int
	LuckyNumbers []int
	DrawnNumbers []int
}

func NewTicket(input string) Ticket {
	var ticket Ticket

	idx1 := strings.IndexRune(input, ':')
	idx2 := strings.IndexRune(input, '|')

	ticket.ID = utils.StringToInt(strings.Trim(input[5:idx1], " "))

	numbers1 := strings.Trim(input[idx1+1:idx2], " ")
	numbers2 := strings.Trim(input[idx2+1:], " ")

	numbers1List := strings.Split(numbers1, " ")
	numbers2List := strings.Split(numbers2, " ")

	for _, number := range numbers1List {
		if number != "" {
			ticket.LuckyNumbers = append(ticket.LuckyNumbers, utils.StringToInt(strings.Trim(number, " ")))
		}
	}

	for _, number := range numbers2List {
		if number != "" {
			ticket.DrawnNumbers = append(ticket.DrawnNumbers, utils.StringToInt(strings.Trim(number, " ")))
		}
	}

	//fmt.Printf("Ticket %d: %v\n", ticket.ID, ticket)
	return ticket
}

func containsValue(list []int, value int) bool {
	for _, num := range list {
		if num == value {
			return true
		}
	}
	return false
}

// ---------------------------------------------------------------------------

func calcWinOfTicket(luckies []int, drawns []int) int {
	var result int

	for _, draw := range drawns {
		if containsValue(luckies, draw) {
			if result == 0 {
				result = 1
			} else {
				result *= 2
			}
		}
	}
	return result
}

func SolvePuzzle1(tickets []Ticket) int {
	var result int = 0

	for _, ticket := range tickets {
		result += calcWinOfTicket(ticket.LuckyNumbers, ticket.DrawnNumbers)
	}

	return result
}

// ---------------------------------------------------------------------------

func calcCopies(luckies []int, drawns []int) int {
	var result int

	for _, draw := range drawns {
		if containsValue(luckies, draw) {
			result++
		}
	}
	return result
}

func SolvePuzzle2(tickets []Ticket) int {
	var result int
	cardCopies := make([]int, len(tickets))

	for idx, ticket := range tickets {
		copies := calcCopies(ticket.LuckyNumbers, ticket.DrawnNumbers)
		cardCopies[idx]++
		if copies > 0 {
			for multiplier := 0; multiplier < cardCopies[idx]; multiplier++ {
				for x := idx + 1; x < len(tickets) && x-idx-1 < copies; x++ {
					cardCopies[x]++
				}
			}
		}
	}

	for _, card := range cardCopies {
		result += card
	}

	return result
}

// ---------------------------------------------------------------------------

func main() {
	var input []string
	var tickets []Ticket

	stopwatch := utils.NewStopwatch()
	stopwatch.Start()

	if runTest {
		input = utils.ReadDataFile(TEST_FILE)
	} else {
		input = utils.ReadDataFile(DATA_FILE)
	}

	for _, line := range input {
		tickets = append(tickets, NewTicket(line))
	}

	result1 := SolvePuzzle1(tickets)
	result2 := SolvePuzzle2(tickets)
	stopwatch.Stop()

	// -------------------------------------
	var elapsedTime time.Duration = stopwatch.GetElapsedTime()
	fmt.Println("Running as test:\t", runTest)
	fmt.Println("Result 1:\t\t\t", result1)
	fmt.Println("Result 2:\t\t\t", result2)
	fmt.Println("Elapsed time:\t\t", elapsedTime)
}
