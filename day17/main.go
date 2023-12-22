// ---------------------------------------------------------------------------
// Golang solution for Advent of Code 2023 day 17
// https://adventofcode.com/2023/day/17
//
// This was created using copilot to assist me in learning Go.
//
// Scenario: Get the crucible from the lava pool to the machine parts factory.
//
//	   To do this, you need to minimize heat loss while choosing a route that
//	   doesn't require the crucible to go in a straight line for too long.
//
//		Input:
//	   lines with number such as "2413432311323"
//	   number = single digit that represents the amount of heat loss
//	            if the crucible enters that block
//
//	 Rule for movement
//	    at most three blocks in a single direction
//	    then turn left or right
//
// Part 1: Start top left, destination bottom right
// Part 2:
// ---------------------------------------------------------------------------
package main

import (
	"container/heap"
	"fmt"
	"math"
	"time"

	"github.com/cdr74/AdventOfCode2023/utils"
)

// runTest defines whether test.data or actual.data should be used
const runTest bool = true
const TEST_FILE string = "test.data"
const DATA_FILE string = "actual.data"

// -------------------------- Common Data Section ----------------------------

const MAX_STRAIGHT = 3

var heatLoss [][]int
var ROWS int
var COLS int

type Position struct {
	r int
	c int
}

type Direction int

const (
	NORTH Direction = iota
	EAST
	SOUTH
	WEST
)

// -------------------------- Queue Shit -------------------------------------

// PriorityQueue implements heap.Interface and holds States.
type PriorityQueue []*State

func (pq PriorityQueue) Len() int { return len(pq) }

func (pq PriorityQueue) Less(i, j int) bool {
	// TODO what's the correct evaluation including streaks?

	if pq[i].heatLoss == pq[j].heatLoss {
		return pq[i].streak < pq[j].streak
	}
	return pq[i].heatLoss < pq[j].heatLoss

	// alternative weight calculation
	//return pq[i].heatLoss+pq[i].streak < pq[j].heatLoss+pq[j].streak
}

func (pq PriorityQueue) Swap(i, j int) { pq[i], pq[j] = pq[j], pq[i] }

func (pq *PriorityQueue) Push(x interface{}) {
	item := x.(*State)
	*pq = append(*pq, item)
}

func (pq *PriorityQueue) Pop() interface{} {
	old := *pq
	n := len(old)
	item := old[n-1]
	*pq = old[0 : n-1]
	return item
}

type StateQueue struct {
	pq PriorityQueue
}

func NewStateQueue() *StateQueue {
	return &StateQueue{pq: make(PriorityQueue, 0)}
}

func (q *StateQueue) Enqueue(state *State) {
	heap.Push(&q.pq, state)
}

func (q *StateQueue) Dequeue() *State {
	if len(q.pq) == 0 {
		return nil
	}
	return heap.Pop(&q.pq).(*State)
}

// -------------------------- Common Code Section ----------------------------

func inputLineToValues(input []string) {
	ROWS = len(input)
	COLS = len(input[0])
	heatLoss = make([][]int, ROWS)
	for i := range heatLoss {
		heatLoss[i] = make([]int, COLS)
	}

	for r, line := range input {
		for c, char := range line {
			heatLoss[r][c] = int(char) - int('0')
		}
	}
}

// -------------------------- Puzzle part 1 ----------------------------------

type State struct {
	position Position
	heatLoss int
	path     []Position
	dir      Direction
	streak   int
}

var accumulatedHeatLoss [][]State

func getNextValidMoves(currentState *State) []Position {
	MAX_STREAK := 3
	var validMoves []Position
	pos := currentState.position

	neighbors := []Position{{pos.r - 1, pos.c}, {pos.r, pos.c + 1}, {pos.r + 1, pos.c}, {pos.r, pos.c - 1}}
	if currentState.dir == NORTH && currentState.streak >= MAX_STREAK {
		neighbors = neighbors[1:] // remove this
	}
	if currentState.dir == EAST && currentState.streak >= MAX_STREAK {
		neighbors = append(neighbors[:1], neighbors[2:]...) // remove this
	}
	if currentState.dir == SOUTH && currentState.streak >= MAX_STREAK {
		neighbors = append(neighbors[:2], neighbors[3:]...) // remove this
	}
	if currentState.dir == WEST && currentState.streak >= MAX_STREAK {
		neighbors = neighbors[:3] // remove this
	}

	for _, neighbor := range neighbors {
		if neighbor.r >= 0 && neighbor.r < ROWS && neighbor.c >= 0 && neighbor.c < COLS {
			if len(currentState.path) > 0 {
				lastPosition := currentState.path[len(currentState.path)-1]
				if neighbor == lastPosition {
					// do not go back to the previous position
					continue
				}
			}
			validMoves = append(validMoves, neighbor)
		}
	}
	return validMoves
}

func getDirection(currentPos Position, nextPos Position) Direction {
	if nextPos.r < currentPos.r {
		return NORTH
	}
	if nextPos.c > currentPos.c {
		return EAST
	}
	if nextPos.r > currentPos.r {
		return SOUTH
	}
	if nextPos.c < currentPos.c {
		return WEST
	}
	return NORTH
}

func exploreBFSStyle(start Position, end Position) int {
	accumulatedHeatLoss = make([][]State, ROWS)
	for i := range accumulatedHeatLoss {
		accumulatedHeatLoss[i] = make([]State, COLS)
	}

	// initialize the queue
	queue := NewStateQueue()
	for r := 0; r < ROWS; r++ {
		for c := 0; c < COLS; c++ {
			state := State{
				position: Position{r, c},
				heatLoss: math.MaxInt32,
				path:     make([]Position, 0),
				dir:      EAST, streak: 0,
			}
			if r == start.r && c == start.c {
				state.heatLoss = 0
			}
			accumulatedHeatLoss[r][c] = state
		}
	}
	queue.Enqueue(&accumulatedHeatLoss[start.r][start.c])

	for queue.pq.Len() > 0 {
		currentState := queue.Dequeue()
		if currentState.heatLoss > accumulatedHeatLoss[currentState.position.r][currentState.position.c].heatLoss {
			// we already found a shorter path to here, skip to next
			continue
		}
		fmt.Printf("Exploring %v --- %v\n", currentState.position, currentState.path)

		if currentState.position == end {
			// TODO: question - stop or not at first solution?
			break
		} else {
			// explore all possible moves
			moveList := getNextValidMoves(currentState)
			currentPos := currentState.position
			for _, move := range moveList {
				// did we find a path with less heat loss to the position indicated by move?
				tmpHeatLoss := accumulatedHeatLoss[currentPos.r][currentPos.c].heatLoss + heatLoss[move.r][move.c]

				// --------------------------------------------------------------------------------------
				// XXX Equal length might be relevant because of streak >> can't overwrite current state
				// --------------------------------------------------------------------------------------
				if tmpHeatLoss <= accumulatedHeatLoss[move.r][move.c].heatLoss {

					direction := getDirection(currentPos, move)
					streak := 1
					if currentState.streak == 0 {
						// first move, no streak (avoid initialization EAST causes streak)
					} else {
						if currentState.dir == direction {
							streak = currentState.streak + 1
						}
					}

					tmpState := State{position: move, heatLoss: tmpHeatLoss, dir: direction, streak: streak}

					tmpState.path = make([]Position, len(currentState.path))
					copy(tmpState.path, currentState.path)
					tmpState.path = append(tmpState.path, move)

					accumulatedHeatLoss[move.r][move.c].path = make([]Position, len(tmpState.path))
					copy(accumulatedHeatLoss[move.r][move.c].path, tmpState.path)

					accumulatedHeatLoss[move.r][move.c].heatLoss = tmpHeatLoss
					fmt.Printf("  Adding %v ---- %d --- %v\n", move, tmpState.heatLoss, tmpState.path)
					queue.Enqueue(&tmpState)
				}
			}
		}
	}

	fmt.Printf("Path: %v\n", accumulatedHeatLoss[end.r][end.c].path)
	return accumulatedHeatLoss[end.r][end.c].heatLoss
}

func SolvePart1() int {
	start := Position{0, 0}
	end := Position{ROWS - 1, COLS - 1}
	cost := exploreBFSStyle(start, end)
	return cost
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

	inputLineToValues(input)

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
