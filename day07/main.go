package main

import (
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/cdr74/AdventOfCode2023/utils"
)

// the flag runTest defines which data file to read
const runTest bool = false
const TEST_FILE string = "test.data"
const DATA_FILE string = "actual.data"

// ---------------------------------------------------------------------------

type Hand struct {
	Cards []int
	Bet   int
}

// MapCardValue maps non-numeric card characters to their corresponding values.
func MapCardValue(char rune) int {
	switch char {
	case 'T':
		return 10
	case 'J':
		return 11
	case 'Q':
		return 12
	case 'K':
		return 13
	case 'A':
		return 14
	default:
		return utils.StringToInt(string(char))
	}
}

// ParseHand parses a string into a Hand struct.
func ParseHand(input string) Hand {
	parts := strings.Fields(input)

	var cards []int
	for _, char := range parts[0] {
		cards = append(cards, MapCardValue(char))
	}

	bet := utils.StringToInt(parts[1])
	return Hand{
		Cards: cards,
		Bet:   bet,
	}
}

func evaluateHand(hand Hand) int {
	cardCount := make(map[int]int)

	for _, card := range hand.Cards {
		cardCount[card]++
	}

	counts := make([]int, 0, len(cardCount))
	for _, count := range cardCount {
		counts = append(counts, count)
	}

	sort.Sort(sort.Reverse(sort.IntSlice(counts)))

	switch {
	case counts[0] == 5:
		return 6
	case counts[0] == 4:
		return 5
	case counts[0] == 3 && counts[1] == 2:
		return 4
	case counts[0] == 3 && counts[1] == 1:
		return 3
	case counts[0] == 2 && counts[1] == 2:
		return 2
	case counts[0] == 2 && counts[1] == 1:
		return 1
	default:
		return 0
	}
}

func rankHands(hand1 Hand, hand2 Hand) int {
	for i := 0; i < len(hand1.Cards); i++ {
		if hand1.Cards[i] > hand2.Cards[i] {
			return 1
		} else if hand1.Cards[i] < hand2.Cards[i] {
			return -1
		}
	}

	return 0
}

// ---------------------------------------------------------------------------

func sortByEvaluation(hands []Hand) {
	sort.Slice(hands, func(i, j int) bool {
		evalI := evaluateHand(hands[i])
		evalJ := evaluateHand(hands[j])

		if evalI == evalJ {
			return rankHands(hands[i], hands[j]) > 0
		}

		return evalI > evalJ
	})
}

func SolvePuzzle1(hands []Hand) int {
	var result int = 0
	sortByEvaluation(hands)
	for idx, hand := range hands {
		fmt.Printf("SolvePuzzle1() - Hand: %v\n", hand)
		result += (len(hands) - idx) * hand.Bet
	}
	return result
}

// ---------------------------------------------------------------------------

func SolvePuzzle2(hands []Hand) int {
	var result int = 0

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

	var hands []Hand
	for _, line := range input {
		hands = append(hands, ParseHand(line))
	}

	result1 := SolvePuzzle1(hands)
	result2 := SolvePuzzle2(hands)
	stopwatch.Stop()

	// -------------------------------------
	var elapsedTime time.Duration = stopwatch.GetElapsedTime()
	fmt.Println("Running as test:\t", runTest)
	fmt.Println("Result 1:\t\t\t", result1)
	fmt.Println("Result 2:\t\t\t", result2)
	fmt.Println("Elapsed time:\t\t", elapsedTime)
}
