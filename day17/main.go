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
// Part 2: same as part 1, different streak length
//
// took inspiration from
// https://www.reddit.com/r/adventofcode/comments/18luw6q/2023_day_17_a_longform_tutorial_on_day_17/
// ---------------------------------------------------------------------------
package main

import (
	"container/heap"
	"fmt"
	"math"

	"github.com/cdr74/AdventOfCode2023/utils"
)

// runTest defines whether test.data or actual.data should be used
const runTest bool = false
const TEST_FILE string = "test.data"
const DATA_FILE string = "actual.data"

var MAX_STREAK = 10 // set to 3 part 1
var MIN_STREAK = 4  // set to 1 part 1

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
	return pq[i].streak < pq[j].streak
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
	dir      Direction
	streak   int
}

func getNextValidMoves(currentState *State) []Position {
	var validMoves []Position
	var neighbors []Position
	pos := currentState.position

	// don't go back, don't exceed streak
	switch currentState.dir {
	case NORTH:
		if currentState.streak < MIN_STREAK {
			// must go north
			neighbors = []Position{{pos.r - 1, pos.c}}
		} else {
			if currentState.streak >= MAX_STREAK {
				// skip north and south
				neighbors = []Position{{pos.r, pos.c + 1}, {pos.r, pos.c - 1}}
			} else {
				// skip south
				neighbors = []Position{{pos.r + 1, pos.c}, {pos.r, pos.c + 1}, {pos.r, pos.c - 1}}
			}
		}
	case EAST:
		if currentState.streak < MIN_STREAK {
			// must go east
			neighbors = []Position{{pos.r, pos.c + 1}}
		} else {
			if currentState.streak >= MAX_STREAK {
				// skip east and west
				neighbors = []Position{{pos.r - 1, pos.c}, {pos.r + 1, pos.c}}
			} else {
				// skip west
				neighbors = []Position{{pos.r - 1, pos.c}, {pos.r, pos.c + 1}, {pos.r + 1, pos.c}}
			}
		}
	case SOUTH:
		if currentState.streak < MIN_STREAK {
			// must go south
			neighbors = []Position{{pos.r + 1, pos.c}}
		} else {
			if currentState.streak >= MAX_STREAK {
				// skip north and south
				neighbors = []Position{{pos.r, pos.c + 1}, {pos.r, pos.c - 1}}
			} else {
				// skip north
				neighbors = []Position{{pos.r, pos.c + 1}, {pos.r + 1, pos.c}, {pos.r, pos.c - 1}}
			}
		}
	case WEST:
		if currentState.streak < MIN_STREAK {
			// must go west
			neighbors = []Position{{pos.r, pos.c - 1}}
		} else {
			if currentState.streak >= MAX_STREAK {
				// skip east and west
				neighbors = []Position{{pos.r - 1, pos.c}, {pos.r + 1, pos.c}}
			} else {
				// skip east
				neighbors = []Position{{pos.r - 1, pos.c}, {pos.r + 1, pos.c}, {pos.r, pos.c - 1}}
			}
		}
	}

	for _, neighbor := range neighbors {
		if neighbor.r >= 0 && neighbor.r < ROWS && neighbor.c >= 0 && neighbor.c < COLS {
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

func lowestCostInStateQueue(stateQueueByCost map[int]*StateQueue) int {
	var lowestCost int = math.MaxInt32
	for cost := range stateQueueByCost {
		if cost < lowestCost && len(stateQueueByCost[cost].pq) > 0 {
			lowestCost = cost
		}
	}
	return lowestCost
}

func findPath(start Position, end Position) int {
	stateQueueByCost := make(map[int]*StateQueue)
	costByStateCache := make(map[State]int)

	startEast := State{position: Position{0, 0}, dir: EAST, streak: 1}
	startSouth := State{position: Position{0, 0}, dir: SOUTH, streak: 1}

	costByStateCache[startEast] = 0
	costByStateCache[startSouth] = 0

	queue := NewStateQueue()
	queue.Enqueue(&startEast)
	queue.Enqueue(&startSouth)

	stateQueueByCost[0] = queue

	for {
		lowestCost := lowestCostInStateQueue(stateQueueByCost)
		queue := stateQueueByCost[lowestCost]

		for len(queue.pq) > 0 {
			currentState := queue.Dequeue()

			if currentState.position == end && currentState.streak >= MIN_STREAK {
				// TODO: don't break yet, test all other paths with same length
				fmt.Printf("Current state: %v\n", currentState)
				fmt.Printf("Heat loss: %v\n", lowestCost)
				return lowestCost
			} else {
				// explore all possible moves
				moveList := getNextValidMoves(currentState)
				for _, move := range moveList {

					// create next state
					direction := getDirection(currentState.position, move)
					streak := 1
					if currentState.dir == direction {
						streak = currentState.streak + 1
					}
					tmpHeatLoss := lowestCost + heatLoss[move.r][move.c]
					tmpState := State{position: move, dir: direction, streak: streak}

					if _, exists := costByStateCache[tmpState]; !exists {
						costByStateCache[tmpState] = tmpHeatLoss

						if _, exists := stateQueueByCost[tmpHeatLoss]; !exists {
							stateQueueByCost[tmpHeatLoss] = NewStateQueue()
						}
						queueByCost := stateQueueByCost[tmpHeatLoss]
						queueByCost.Enqueue(&tmpState)
					}
				}
			}
		}
	}
}

func SolvePart1() int {
	start := Position{0, 0}
	end := Position{ROWS - 1, COLS - 1}
	cost := findPath(start, end)
	return cost
}

// -------------------------- Main entry -------------------------------------

func main() {
	var input []string

	if runTest {
		input = utils.ReadDataFile(TEST_FILE)
	} else {
		input = utils.ReadDataFile(DATA_FILE)
	}

	inputLineToValues(input)
	result1 := SolvePart1()

	// ---------------------- Print results ----------------------------------
	fmt.Println("Running as test:\t", runTest)
	fmt.Println("Result 1:\t\t", result1)
}
