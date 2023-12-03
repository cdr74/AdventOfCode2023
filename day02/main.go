package main

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/cdr74/AdventOfCode2023/utils"
)

// the flag runTest defines which data file to read
const runTest bool = false
const TEST_FILE string = "test.data"
const DATA_FILE string = "actual.data"

// ---------------------------------------------------------------------------

type Bag struct {
	RedItems   int
	BlueItems  int
	GreenItems int
}

type Game struct {
	ID   int
	bags []Bag
}

// Input looks like "6 red, 1 blue, 3 green" each color is optional, set 0 if not present
func stringToBag(input string) Bag {
	var bag Bag = Bag{RedItems: 0, BlueItems: 0, GreenItems: 0}

	colors := strings.Split(input, ",")
	for _, color := range colors {
		color = strings.Trim(color, " ")
		parts := strings.Split(color, " ")
		value, _ := strconv.Atoi(parts[0])
		switch parts[1] {
		case "red":
			bag.RedItems = value
		case "blue":
			bag.BlueItems = value
		case "green":
			bag.GreenItems = value
		default:
			panic("Color not defined [" + parts[1] + "]")
		}
	}
	return bag
}

// Turns an input line into a Game type. Each input line looks like
// Game 5: 6 red, 1 blue, 3 green; 2 blue, 1 red, 2 green
// There's always a prefix of "Game x:" Color sets are seperated by;
// Any color is optional
func stringToToGame(input string) Game {
	var game Game

	spaceIndex := strings.Index(input, " ")
	colonIndex := strings.Index(input, ":")
	numberSubstring := input[spaceIndex+1 : colonIndex]
	game.ID, _ = strconv.Atoi(numberSubstring)

	// ignore first part with game id, read colors
	input = input[colonIndex+1:]

	// iterate over all bag draws in a game
	gameStrings := strings.Split(input, ";")
	for _, gameString := range gameStrings {
		bag := stringToBag(gameString)
		game.bags = append(game.bags, bag)
	}
	return game
}

// Iterate over input file represented by list of strings.
// returns a list of Games
func inputToGames(input []string) []Game {
	var games []Game
	for _, line := range input {
		gamesOfLine := stringToToGame(line)
		games = append(games, gamesOfLine)
	}

	return games
}

// ---------------------------------------------------------------------------

func isGameValid(game Game, bag Bag) bool {
	for _, gameBag := range game.bags {
		if gameBag.BlueItems > bag.BlueItems ||
			gameBag.GreenItems > bag.GreenItems ||
			gameBag.RedItems > bag.RedItems {
			return false
		}
	}
	return true
}

// Finds the ID's of possible games (bag colors are in range) and returns the summ of ID's
func SolvePuzzle1(games []Game) int {
	var bagContentGame1 Bag = Bag{RedItems: 12, BlueItems: 14, GreenItems: 13}
	var result int = 0
	for _, game := range games {
		if isGameValid(game, bagContentGame1) {
			result += game.ID
		}
	}
	return result
}

// ---------------------------------------------------------------------------

// initialize with one as power calc fails otherwise
func getMinimalBag(game Game) Bag {
	var minimalBag Bag = Bag{RedItems: 1, BlueItems: 1, GreenItems: 1}
	for _, bag := range game.bags {
		if bag.BlueItems > minimalBag.BlueItems {
			minimalBag.BlueItems = bag.BlueItems
		}
		if bag.RedItems > minimalBag.RedItems {
			minimalBag.RedItems = bag.RedItems
		}
		if bag.GreenItems > minimalBag.GreenItems {
			minimalBag.GreenItems = bag.GreenItems
		}
	}
	return minimalBag
}

// Finds for each game the minimal bag, calculates the power of the minima bag
// power = red * blue * green, then returns the sum of power
func SolvePuzzle2(games []Game) int {
	var result int = 0
	for _, game := range games {
		minimalBag := getMinimalBag(game)
		power := minimalBag.RedItems * minimalBag.BlueItems * minimalBag.GreenItems
		result += power
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

	games := inputToGames(input)
	result1 := SolvePuzzle1(games)
	result2 := SolvePuzzle2(games)
	stopwatch.Stop()

	// -------------------------------------
	var elapsedTime time.Duration = stopwatch.GetElapsedTime()
	fmt.Println("Running as test:\t", runTest)
	fmt.Println("Result 1:\t\t\t", result1)
	fmt.Println("Result 2:\t\t\t", result2)
	fmt.Println("Elapsed time:\t\t", elapsedTime)
}
