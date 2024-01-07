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
// ground is at z=0, so lowest is 1
//
// Part 1:
//
//	let brick fall down until they touch ground or a brick below them
//	identify which bricks can be disintegrated (do not support other bricks)
//
// Part 2:
//
//	sum up for all bricks that are structurally relevant, howm many
//	other bricks falls down
//
// ---------------------------------------------------------------------------
package main

import (
	"fmt"
	"sort"
	"strings"

	"github.com/cdr74/AdventOfCode2023/utils"
	"github.com/pkg/profile"
	"github.com/samber/lo"
)

// runTest defines whether test.data or actual.data should be used
const runTest bool = false
const TEST_FILE string = "test.data"
const DATA_FILE string = "actual.data"

// -------------------------- Common Code Section ----------------------------

var brickMapByLowestZ map[int][]*Brick
var brickCount int = 0

func deepCopyBrickMapByLowestZ(org map[int][]*Brick) map[int][]*Brick {
	newMap := make(map[int][]*Brick)
	for lowestZ, bricks := range org {
		newBricks := make([]*Brick, len(bricks))
		for i, brick := range bricks {
			newBricks[i] = brick.Copy()
		}
		newMap[lowestZ] = newBricks
	}
	return newMap
}

func removeBrickFromBrickMap(ID int) {
	for i := range brickMapByLowestZ {
		brickMapByLowestZ[i] = lo.Filter(brickMapByLowestZ[i], func(b *Brick, i int) bool {
			return b.ID != ID
		})
	}
}

// Brick represents a 3D brick
type Brick struct {
	ID                     int
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

func (b *Brick) intersectionFree() bool {
	otherBricks := brickSliceFromMap()
	for _, otherBrick := range otherBricks {
		if b.ID != otherBrick.ID && b.Intersects(otherBrick) {
			return false
		}
	}
	return true
}

func (b *Brick) Copy() *Brick {
	return &Brick{
		ID:     b.ID,
		StartX: b.StartX,
		StartY: b.StartY,
		StartZ: b.StartZ,
		EndX:   b.EndX,
		EndY:   b.EndY,
		EndZ:   b.EndZ,
	}
}

// CreateBricks creates a list of bricks
func CreateBricks(input []string) {
	brickMapByLowestZ = make(map[int][]*Brick)

	for _, line := range input {
		coordString := strings.Replace(line, "~", ",", -1)
		points := strings.Split(coordString, ",")
		if len(points) != 6 {
			panic("CreateBricks() - Invalid number of points in line")
		}

		brick := Brick{
			ID:     brickCount,
			StartX: utils.StringToInt(points[0]),
			StartY: utils.StringToInt(points[1]),
			StartZ: utils.StringToInt(points[2]),
			EndX:   utils.StringToInt(points[3]),
			EndY:   utils.StringToInt(points[4]),
			EndZ:   utils.StringToInt(points[5]),
		}
		brickCount++

		if brickMapByLowestZ[brick.StartZ] == nil {
			brickMapByLowestZ[brick.StartZ] = make([]*Brick, 0)
		}
		brickMapByLowestZ[brick.StartZ] = append(brickMapByLowestZ[brick.StartZ], &brick)
	}
}

func brickSliceFromMap() []*Brick {
	bricks := make([]*Brick, 0)
	for i := range brickMapByLowestZ {
		lo.ForEach(brickMapByLowestZ[i], func(b *Brick, _ int) {
			bricks = append(bricks, b.Copy())
		})
	}
	return bricks
}

// -------------------------- Puzzle part 1 ----------------------------------

func moveBrickDownByOne(brick *Brick) bool {
	if brick.StartZ <= 1 {
		return false
	}
	// update brick coordinates
	brick.EndZ--
	brick.StartZ--

	// check if brick intersects with any brick
	movePossible := brick.intersectionFree()

	if movePossible {
		// move brick to new level in brickMapByLowestZ
		if brickMapByLowestZ[brick.StartZ] == nil {
			brickMapByLowestZ[brick.StartZ] = make([]*Brick, 0)
		}
		// add brick to new level
		brickMapByLowestZ[brick.StartZ] = append(brickMapByLowestZ[brick.StartZ], brick)
		// remove from old level in brickMapByLowestZ
		oldZ := brick.StartZ + 1
		brickMapByLowestZ[oldZ] = lo.Filter(brickMapByLowestZ[oldZ], func(b *Brick, i int) bool {
			return *b != *(brick)
		})
	} else {
		// undo move
		brick.EndZ++
		brick.StartZ++
		return false
	}

	return true
}

func moveBricksAtZDownByOne(z int, fallenBricks []bool) int {
	moveCnt := 0

	for _, brick := range brickMapByLowestZ[z] {
		moving := true
		for moving {
			if moveBrickDownByOne(brick) {
				fallenBricks[brick.ID] = true
				moveCnt++
			} else {
				moving = false
			}
		}
	}

	return moveCnt
}

func letBricksFallDown() int {
	fallenBricks := make([]bool, brickCount)

	// let bricks fall down until they touch ground or a brick below them
	done := false
	for !done {
		done = true
		zKeys := lo.Keys(brickMapByLowestZ)
		sort.Ints(zKeys)

		for _, z := range zKeys {
			// try to move each brick down if not at 1 (lowest)
			if z <= 1 {
				continue
			}

			moved := moveBricksAtZDownByOne(z, fallenBricks)
			if moved > 0 {
				done = false
			}
		}
	}

	total := 0
	for f := range fallenBricks {
		if fallenBricks[f] {
			total++
		}
	}
	return total
}

func countStructuralRelevantBricks() (int, int) {
	totalMovedCount := 0
	structuralRelevantBricks := 0

	copy_brickMapByLowestZ := deepCopyBrickMapByLowestZ(brickMapByLowestZ)
	loop_bricks := brickSliceFromMap()

	for _, brick := range loop_bricks {
		// remove this brick, see if anything would move
		removeBrickFromBrickMap(brick.ID)

		bricksMovedCount := letBricksFallDown()
		if bricksMovedCount > 0 {
			totalMovedCount += bricksMovedCount
			structuralRelevantBricks++
			fmt.Printf("%v --v-- %d\n", brick, bricksMovedCount)
		}

		// restore state of global brick data structure
		brickMapByLowestZ = deepCopyBrickMapByLowestZ(copy_brickMapByLowestZ)
	}

	return len(loop_bricks) - structuralRelevantBricks, totalMovedCount
}

func SolvePart1_2() (int, int) {
	letBricksFallDown()
	// find all bricks that are needed to support other bricks
	b, t := countStructuralRelevantBricks()
	return b, t
}

// -------------------------- Main entry -------------------------------------

func main() {
	// Enable CPU profiling
	defer profile.Start(profile.ProfilePath(".")).Stop()

	var input []string

	if runTest {
		input = utils.ReadDataFile(TEST_FILE)
	} else {
		input = utils.ReadDataFile(DATA_FILE)
	}

	CreateBricks(input)
	result1, result2 := SolvePart1_2()

	// ---------------------- Print results ----------------------------------

	fmt.Println("Running as test:\t", runTest)
	fmt.Println("Result 1:\t\t", result1)
	fmt.Println("Result 2:\t\t", result2)
}
