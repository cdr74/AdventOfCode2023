// ---------------------------------------------------------------------------
// Golang solution for Advent of Code 2023 day ..
// https://adventofcode.com/2023/day/..
//
// This was created using copilot to assist me in learning Go.
//
// Scenario:
//
//		Read a list of 2d vectors from a file. Dat in the file has 3 diemnsions,
//		but the 3rd dimension z is to be ignored.
//
//		Data format:
//	    19, 13, 30 @ -2,  1, -2 (x, y, z, @, dx, dy, dz)
//
// Part 1:
// Part 2:
// ---------------------------------------------------------------------------
package main

import (
	"fmt"
	"math"
	"strings"
	"time"

	"github.com/cdr74/AdventOfCode2023/utils"
)

// runTest defines whether test.data or actual.data should be used
const runTest bool = false
const TEST_FILE string = "test.data"
const DATA_FILE string = "actual.data"

// -------------------------- Common Code Section ----------------------------

type Point struct {
	X float64
	Y float64
}

type Vector2D struct {
	X  float64
	Y  float64
	DX float64
	DY float64
}

func stringToVector2D(s string) Vector2D {
	var x, y, z, dx, dy, dz float64

	_, err := fmt.Sscanf(strings.ReplaceAll(s, ",", ""), "%f %f %f @ %f %f %f", &x, &y, &z, &dx, &dy, &dz)
	if err != nil {
		panic("Error parsing vector data")
	}

	v := Vector2D{X: x, Y: y, DX: dx, DY: dy}
	return v
}

func inputToVectorList(input []string) []Vector2D {
	var vectors []Vector2D

	for _, s := range input {
		v := stringToVector2D(s)
		vectors = append(vectors, v)
	}
	return vectors
}

// -------------------------- Puzzle part 1 ----------------------------------

/*
 *  Over n ticks of time t the vector will move to (x,y)
 *  (x,y) = (Vector2D.DX, Vector2D.DY) * t + (Vector2D.X , Vector2D.Y)
 *
 *  to find the intersection set both equal same (x, y)  and solve for t1, t2
 */
func intersectionPoint(v1 Vector2D, v2 Vector2D) Point {
	var p = Point{X: 0, Y: 0}

	t1 := (v2.Y - v1.Y + (v2.DY/v2.DX)*v1.X - (v2.DY/v2.DX)*v2.X) / (v1.DY - (v2.DY/v2.DX)*v1.DX)
	t2 := (v1.X + v1.DX*t1 - v2.X) / v2.DX

	if math.IsInf(t1, 0) || math.IsNaN(t1) || t1 < float64(0) {
		return p
	}
	if math.IsInf(t2, 0) || math.IsNaN(t2) || t2 < float64(0) {
		return p
	}

	p.X = v1.DX*t1 + v1.X
	p.Y = v1.DY*t1 + v1.Y

	fmt.Printf("OK: v1: %v, v2: %v, t1: %f, p2: %v\n", v1, v2, t1, p)
	return p
}

func SolvePart1(vectors []Vector2D) int {
	minDistance := 7
	maxDistance := 27

	if !runTest {
		minDistance = 200000000000000
		maxDistance = 400000000000000
	}

	count := 0

	for x := 0; x < len(vectors)-1; x++ {
		v1 := vectors[x]
		for y := x + 1; y < len(vectors); y++ {
			v2 := vectors[y]

			p := intersectionPoint(v1, v2)

			if p.X >= float64(minDistance) && p.X <= float64(maxDistance) &&
				p.Y >= float64(minDistance) && p.Y <= float64(maxDistance) {
				count++
			}
		}
	}
	return count
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

	vectors := inputToVectorList(input)
	result1 := SolvePart1(vectors)
	result2 := SolvePart2()
	stopwatch.Stop()

	// ---------------------- Print results ----------------------------------
	var elapsedTime time.Duration = stopwatch.GetElapsedTime()
	fmt.Println("Running as test:\t", runTest)
	fmt.Println("Result 1:\t\t", result1)
	fmt.Println("Result 2:\t\t", result2)
	fmt.Println("Elapsed time:\t\t", elapsedTime)
}
