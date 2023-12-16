// ---------------------------------------------------------------------------
// Golang solution for Advent of Code 2023 day 16
// https://adventofcode.com/2023/day/16
//
// This was created using copilot to assist me in learning Go.
//
// Scenario: 2 dimensional array of tiles, a laser enters from top left.
// The laser moves in a straight line until it hits a mirror or a splitter.
//
// Laser hitting mitror: Laser is reflected 90 degrees
// Laser hitting splitter: if hitting the pointy end, nothing happens, if hitting
// the flat end, the laser is split and sent in the directions of both pointy ends.
//
// If a beam passes through a tile it is energized (empty tiles, mirrors and splitters).
//
// Valid characters representing the field are:
//
//	. = floor
//	|,- = splitter
//	/,\ = mirror
//
// Part 1: count the number of energized tiles after the laser has passed through
// Part 2: find the start position that leads to the highest number of energized tiles
// ---------------------------------------------------------------------------
package main

import (
	"fmt"
	"hash/fnv"
	"time"

	"github.com/cdr74/AdventOfCode2023/utils"
)

// runTest defines whether test.data or actual.data should be used
const runTest bool = false
const TEST_FILE string = "test.data"
const DATA_FILE string = "actual.data"

// -------------------------- Common Data Section ----------------------------

type Tile struct {
	IsEnergized          bool
	IsSplitterNorthSouth bool
	IsSplitterEastWest   bool
	IsMirrorUpRight      bool
	IsMirrorDownRight    bool
	IsFloor              bool
}

type Direction int

const (
	NORTH Direction = iota
	EAST
	WEST
	SOUTH
)

type Laser struct {
	Row       int
	Col       int
	direction Direction
}

var cache = make(map[uint32]int)
var beams []Laser
var field [][]Tile
var ROWS int
var COLS int

// -------------------------- Common Code Section ----------------------------

func HashBeam(beam Laser) uint32 {
	hash := fnv.New32()
	hash.Write([]byte(fmt.Sprintf("%v", beam)))
	return hash.Sum32()
}

func inputLineToValues(input []string) {
	ROWS = len(input)
	COLS = len(input[0])
	field = make([][]Tile, ROWS)
	for i := range field {
		field[i] = make([]Tile, COLS)
	}

	for r, line := range input {
		for c, char := range line {
			tile := Tile{IsEnergized: false, IsSplitterNorthSouth: false, IsSplitterEastWest: false, IsMirrorUpRight: false, IsMirrorDownRight: false, IsFloor: true}
			switch char {
			case '.':
				tile.IsFloor = true
			case '/':
				tile.IsMirrorUpRight = true
			case '\\':
				tile.IsMirrorDownRight = true
			case '-':
				tile.IsSplitterEastWest = true
			case '|':
				tile.IsSplitterNorthSouth = true
			}
			field[r][c] = tile
		}
	}
}

func printField() {
	for r := 0; r < ROWS; r++ {
		for c := 0; c < COLS; c++ {
			tile := field[r][c]
			if tile.IsEnergized {
				fmt.Print("X")
			} else {
				fmt.Print(".")
			}
		}
		fmt.Println()
	}
	fmt.Printf("\n\n")
}

func addBeam(row int, col int, direction Direction) {
	beam := Laser{Row: row, Col: col, direction: direction}
	beamHash := HashBeam(beam)
	if _, ok := cache[beamHash]; ok {
		// beam has been seen before
		return
	}
	cache[beamHash] = 1
	beams = append(beams, beam)

}

// -------------------------- Puzzle part 1 ----------------------------------

// todo: cover loops
func runBeamTillEnd(beam Laser) {
	var positionCache = make(map[Laser]bool)
	notInLoop := true
	for notInLoop {
		tile := &field[beam.Row][beam.Col]
		tile.IsEnergized = true

		if tile.IsMirrorDownRight {
			switch beam.direction {
			case NORTH:
				beam.direction = WEST
			case EAST:
				beam.direction = SOUTH
			case WEST:
				beam.direction = NORTH
			case SOUTH:
				beam.direction = EAST
			}
		}

		if tile.IsMirrorUpRight {
			switch beam.direction {
			case NORTH:
				beam.direction = EAST
			case EAST:
				beam.direction = NORTH
			case WEST:
				beam.direction = SOUTH
			case SOUTH:
				beam.direction = WEST
			}
		}

		// create a new beam if we're north / south headed and put it into beams
		// else treat as normal field and pass through
		if tile.IsSplitterEastWest && (beam.direction == NORTH || beam.direction == SOUTH) {
			beam.direction = EAST
			if beam.Col > 0 {
				addBeam(beam.Row, beam.Col-1, WEST)
			}
		}

		// create a new beam if we're east / west headed and put it into beams
		// else treat as normal field and pass through
		if tile.IsSplitterNorthSouth && (beam.direction == EAST || beam.direction == WEST) {
			beam.direction = NORTH
			if beam.Row < ROWS-1 {
				addBeam(beam.Row+1, beam.Col, SOUTH)
			}
		}

		switch beam.direction {
		case NORTH:
			beam.Row--
		case EAST:
			beam.Col++
		case WEST:
			beam.Col--
		case SOUTH:
			beam.Row++
		}
		if _, ok := positionCache[beam]; ok {
			// beam has been seen before
			notInLoop = false
		}
		positionCache[Laser{Row: beam.Row, Col: beam.Col, direction: beam.direction}] = true

		if beam.Row < 0 || beam.Row >= ROWS || beam.Col < 0 || beam.Col >= COLS {
			// beam has left the field
			return
		}
	}
}

func countEnergy() int {
	var result int = 0
	for r := 0; r < ROWS; r++ {
		for c := 0; c < COLS; c++ {
			if field[r][c].IsEnergized {
				result++
			}
		}
	}
	return result
}

func SolvePart1() int {
	beam := Laser{Row: 0, Col: 0, direction: EAST}
	beams = append(beams, beam)

	for len(beams) > 0 {
		beam := beams[0]
		beams = beams[1:]
		runBeamTillEnd(beam)
	}

	printField()
	return countEnergy()
}

// -------------------------- Puzzle part 2 ----------------------------------

func solveWithStart(beam Laser, input []string) int {
	inputLineToValues(input)
	cache = make(map[uint32]int)
	beams = append(beams, beam)
	for len(beams) > 0 {
		beam := beams[0]
		beams = beams[1:]
		runBeamTillEnd(beam)
	}

	return countEnergy()
}

func SolvePart2(input []string) int {
	var result int = 0
	// left to right
	for r := 0; r < ROWS; r++ {
		beam := Laser{Row: r, Col: 0, direction: EAST}
		count := solveWithStart(beam, input)
		if count > result {
			result = count
		}
	}
	// right to left
	for r := 0; r < ROWS; r++ {
		beam := Laser{Row: r, Col: COLS - 1, direction: WEST}
		count := solveWithStart(beam, input)
		if count > result {
			result = count
		}
	}
	// top to bottom
	for c := 0; c < COLS; c++ {
		beam := Laser{Row: 0, Col: c, direction: SOUTH}
		count := solveWithStart(beam, input)
		if count > result {
			result = count
		}
	}
	// bottom to top
	for c := 0; c < COLS; c++ {
		beam := Laser{Row: ROWS - 1, Col: c, direction: NORTH}
		count := solveWithStart(beam, input)
		if count > result {
			result = count
		}
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

	inputLineToValues(input)
	result1 := SolvePart1()
	result2 := SolvePart2(input)

	stopwatch.Stop()

	// ---------------------- Print results ----------------------------------
	var elapsedTime time.Duration = stopwatch.GetElapsedTime()
	fmt.Println("Running as test:\t", runTest)
	fmt.Println("Result 1:\t\t", result1)
	fmt.Println("Result 2:\t\t", result2)
	fmt.Println("Elapsed time:\t\t", elapsedTime)
}
