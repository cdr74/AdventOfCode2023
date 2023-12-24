// ---------------------------------------------------------------------------
// Golang solution for Advent of Code 2023 day 22
// https://adventofcode.com/2023/day/22
//
// This was created using copilot to assist me in learning Go.
//
// Scenario:

// Data
// Each line of text represents the position of a single brick at a time
// The position is given as two x,y,z coordinates
// - one for each end of the brick
// - separated by a tilde (~)
// Each brick is made up of a single straight line of cubes
// The whole snapshot is aligned to a three-dimensional cube grid.
//
// A line like 2,2,2~2,2,2 means that both ends of the brick are at the same coordinate
// in other words, that the brick is a single cube.
//
// ground is at z=0
//
// Part 1: let brick fall down until they touch ground or a brick below them
//
//	identify which bricks can be disintegrated (do not support other bricks)
//
// Part 2:
// ---------------------------------------------------------------------------
package main

import (
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/cdr74/AdventOfCode2023/utils"
)

// runTest defines whether test.data or actual.data should be used
const runTest bool = true
const TEST_FILE string = "test.data"
const DATA_FILE string = "actual.data"

// -------------------------- Common Code Section ----------------------------

var brickMapByLowestZ map[int][]*Brick
var bricks []*Brick

func removeBrickFromBrickMap(brick *Brick) {
	slice := brickMapByLowestZ[brick.StartZ]
	s := make([]*Brick, 0)
	for idx, _ := range slice {
		if *(slice[idx]) != *(brick) {
			s = append(s, slice[idx])
		}
	}
	brickMapByLowestZ[brick.StartZ] = s
}

func removeBrickFromList(brick *Brick) {
	s := make([]*Brick, 0)
	for idx, _ := range bricks {
		if *(bricks[idx]) != *(brick) {
			s = append(s, bricks[idx])
		}
	}
	bricks = s
}

// Brick represents a 3D brick
type Brick struct {
	StartX, StartY, StartZ int
	EndX, EndY, EndZ       int
}

// Intersects checks if two bricks intersect
func (b *Brick) Intersects(other *Brick) bool {
	return b.StartX <= other.EndX &&
		b.EndX >= other.StartX &&
		b.StartY <= other.EndY &&
		b.EndY >= other.StartY &&
		b.StartZ <= other.EndZ &&
		b.EndZ >= other.StartZ
}

// CreateBricks creates a list of bricks
func CreateBricks(input []string) {
	brickMapByLowestZ = make(map[int][]*Brick)
	bricks = make([](*Brick), 0)
	for _, line := range input {
		coordString := strings.Replace(line, "~", ",", -1)
		points := strings.Split(coordString, ",")
		if len(points) != 6 {
			panic("CreateBricks() - Invalid number of points in line")
		}

		startX := utils.StringToInt(points[0])
		startY := utils.StringToInt(points[1])
		startZ := utils.StringToInt(points[2])
		endX := utils.StringToInt(points[3])
		endY := utils.StringToInt(points[4])
		endZ := utils.StringToInt(points[5])

		if startX > endX || startY > endY || startZ > endZ {
			panic("CreateBricks() - start coordinates are not smaller than end coordinates")
		}

		brick := Brick{
			StartX: startX,
			StartY: startY,
			StartZ: startZ,
			EndX:   endX,
			EndY:   endY,
			EndZ:   endZ,
		}

		bricks = append(bricks, &brick)
		if brickMapByLowestZ[brick.StartZ] == nil {
			brickMapByLowestZ[brick.StartZ] = make([]*Brick, 0)
		}
		brickMapByLowestZ[brick.StartZ] = append(brickMapByLowestZ[brick.StartZ], &brick)
	}
}

// -------------------------- Puzzle part 1 ----------------------------------

func moveBricksAtZDownByOne(z int) int {
	moveCnt := 0

	for _, brick := range brickMapByLowestZ[z] {
		if brick.StartZ <= 1 {
			continue
		}

		// update brick coordinates
		brick.EndZ--
		brick.StartZ--
		// check if brick intersects with any brick
		movePossible := true
		for _, otherBrick := range bricks {
			if !(otherBrick == brick) && brick.Intersects(otherBrick) {
				movePossible = false
				break
			}
		}
		if movePossible {
			// move to new level in brickMapByLowestZ
			moveCnt++
			if brickMapByLowestZ[brick.StartZ] == nil {
				brickMapByLowestZ[brick.StartZ] = make([]*Brick, 0)
			}
			// add brick to new level
			brickMapByLowestZ[brick.StartZ] = append(brickMapByLowestZ[brick.StartZ], brick)
			// remove from old level in brickMapByLowestZ
			removeBrickFromBrickMap(brick)
		} else {
			// undo move
			brick.EndZ++
			brick.StartZ++
		}

	}
	return moveCnt
}

func letBricksFallDown() int {
	// let bricks fall down until they touch ground or a brick below them
	bricksMovedCount := 0
	loopCnt := 0
	for loopCnt > 0 {
		loopCnt = 0

		// get all keys from brickMapByLowestZ
		zKeys := make([]int, 0)
		for k := range brickMapByLowestZ {
			if brickMapByLowestZ[k] == nil || len(brickMapByLowestZ[k]) == 0 {
				continue
			}
			zKeys = append(zKeys, k)
		}
		// sort keys
		sort.Ints(zKeys)

		for _, z := range zKeys {
			// try to move each brick down if not at 1 (lowest)
			if z <= 1 {
				continue
			}
			bricksMovedCount += moveBricksAtZDownByOne(z)
			loopCnt++
		}
	}
	return bricksMovedCount
}

func deepCopyOfBricks() []*Brick {
	copy_bricks := make([](*Brick), 0)
	for _, b := range bricks {
		brick := Brick{b.StartX, b.StartY, b.StartZ, b.EndX, b.EndY, b.EndZ}
		copy_bricks = append(copy_bricks, &brick)
	}
	return copy_bricks
}

func deepCopyOfMapByLowestZ() map[int][]*Brick {
	copy_brickMapByLowestZ := make(map[int][]*Brick)
	for k, slice := range brickMapByLowestZ {
		copy_slice := make([]*Brick, 0)
		for _, b := range slice {
			brick := Brick{b.StartX, b.StartY, b.StartZ, b.EndX, b.EndY, b.EndZ}
			copy_slice = append(copy_slice, &brick)
		}
		copy_brickMapByLowestZ[k] = copy_slice
	}
	return copy_brickMapByLowestZ
}

func countStructuralRelevantBricks() int {
	var copy_brickMapByLowestZ map[int][]*Brick
	var copy_bricks []*Brick
	structuralRelevantBricks := 0
	var loop_bricks = deepCopyOfBricks()

	for _, brick := range loop_bricks {
		// save state of bricks and brickMapByLowestZ
		copy_brickMapByLowestZ = deepCopyOfMapByLowestZ()
		copy_bricks = deepCopyOfBricks()

		// remove this brick, see if anything would move
		removeBrickFromBrickMap(brick)
		// remove from bricks
		removeBrickFromList(brick)

		bricksMovedCount := letBricksFallDown()
		if bricksMovedCount > 0 {
			structuralRelevantBricks++
		}

		// restore state of global structures
		brickMapByLowestZ = copy_brickMapByLowestZ
		bricks = copy_bricks
	}

	return structuralRelevantBricks
}

func SolvePart1() int {
	letBricksFallDown()
	for _, brick := range bricks {
		fmt.Printf("%v\n", brick)
	}
	// find all bricks that are needed to support other bricks
	return countStructuralRelevantBricks()
}

// -------------------------- Puzzle part 2 ----------------------------------

func SolvePart2() int {
	var result int = 0

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

	CreateBricks(input)
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
