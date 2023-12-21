// ---------------------------------------------------------------------------
// Golang solution for Advent of Code 2023 day 20
// https://adventofcode.com/2023/day/20
//
// This was created using copilot to assist me in learning Go.
//
// Scenario:
//
//		Data describes a network with 3 types of nodes that are connected to one another.
//
//	 Node types:
//	   broadcaster -> a, b  = start of signal distribution can have 1..n nodes connected
//	   %gf -> fn, lx        = flip flop node, can have 1..n nodes connected
//	   &db -> kg, sp		   = Conjunction nonde, can have 1..2 nodes connected
//
//		  Not every node has a follow up (there's potential end nodes, not marked as such)
//
// Part 1: issue 1000 low pulses, result is the number of low pulses * number of high pulses
// Part 2:
// ---------------------------------------------------------------------------
package main

import (
	"fmt"
	"strings"
	"time"

	"github.com/cdr74/AdventOfCode2023/utils"
)

// runTest defines whether test.data or actual.data should be used
const runTest bool = false
const TEST_FILE string = "test.data"
const DATA_FILE string = "actual.data"

// -------------------------- Common Code Section ----------------------------

var done = false
var buttonCount = 0
var lowPulseCount = 0
var highPulseCount = 0

var nodeStore = make(map[string]*BasicNode)

// Pulse represents a digital signal pulse.
type Pulse int

// Constants for high and low pulses.
const (
	Low Pulse = iota
	High
)

// Pulse represents a digital signal pulse.
type NodeType int

// Constants for high and low pulses.
const (
	Basic NodeType = iota
	FlipFlop
	Conjunction
)

type Node interface {
	receivePulse(fromNode string, pulse Pulse)
}

type BasicNode struct {
	name         string
	destinations []string
	nodeType     NodeType
	isOn         bool
	memory       map[string]Pulse
}

// just forwards Pulse to all connected nodes
func (b *BasicNode) receivePulse(fromNode string, pulse Pulse) {
	switch b.nodeType {
	case Basic:
		// just forward pulse to all destinations
		for _, n := range b.destinations {
			messageQueue.Push(b, n, pulse)
		}
	case FlipFlop:
		b.flipFlopReceivePulse(fromNode, pulse)
	case Conjunction:
		b.conjunctionReceivePulse(fromNode, pulse)
	}
}

// -------------------------------------------
// isOn -> !isOn
// Low pulse
// - if it was off -> send High Pulse
// - if it was on -> send Low Pulse
//
// High Pulse
// - any isOn -> ignore
func (b *BasicNode) flipFlopReceivePulse(fromNode string, pulse Pulse) {
	if pulse == High {
		return
	}
	var newPulse Pulse

	if b.isOn {
		b.isOn = false
		newPulse = Low
	} else {
		b.isOn = true
		newPulse = High
	}

	for _, n := range b.destinations {
		// put pulse as message to all destinations in a queue
		messageQueue.Push(b, n, newPulse)
	}
}

func (b *BasicNode) initializeMemory(origins []string) {
	b.memory = make(map[string]Pulse)
	for _, o := range origins {
		b.memory[o] = Low
	}
}

func (c *BasicNode) allMemoryHigh() bool {
	for _, v := range c.memory {
		if v == Low {
			return false
		}
	}
	return true
}

// -------------------------------------------
// Low pulse:
// - send high pulse
// High Pulse
// - if memory has only high pulses -> send Low pulse
// - otherwise send high pulse
func (b *BasicNode) conjunctionReceivePulse(fromNode string, pulse Pulse) {
	// update memory
	b.memory[fromNode] = pulse
	var newPulse Pulse
	if b.allMemoryHigh() {
		// send Low pulse to all destinations
		newPulse = Low
	} else {
		// send High pulse to all destinations
		newPulse = High
	}

	for _, n := range b.destinations {
		// put pulse as message to all destinations in a queue
		messageQueue.Push(b, n, newPulse)
	}
}

// ---------------------------------------------------------------------------

type Message struct {
	from  *BasicNode
	to    string
	pulse Pulse
}

type MessageQueue struct {
	messages []Message
}

func NewMessageQueue() *MessageQueue {
	return &MessageQueue{
		messages: make([]Message, 0),
	}
}

func (mq *MessageQueue) Push(from *BasicNode, to string, pulse Pulse) {
	message := Message{
		from:  from,
		to:    to,
		pulse: pulse,
	}
	if pulse == Low {
		lowPulseCount++
	} else {
		highPulseCount++
	}

	mq.messages = append(mq.messages, message)
}

func (mq *MessageQueue) Pop() Message {
	if len(mq.messages) == 0 {
		// Returning a zero-value Message when the queue is empty
		return Message{}
	}

	message := mq.messages[0]
	mq.messages = mq.messages[1:]
	return message
}

var messageQueue = NewMessageQueue()

// ---------------------------------------------------------------------------

func createNode(nodeName string, nodeType string) *BasicNode {
	var node BasicNode
	if nodeType == "%" {
		// create flip flop node
		node = BasicNode{name: nodeName, nodeType: FlipFlop, isOn: false}
	} else if nodeType == "&" {
		// create conjunction node
		node = BasicNode{name: nodeName, nodeType: Conjunction, isOn: false}
	} else {
		// create basic node
		node = BasicNode{name: nodeName, nodeType: Basic, isOn: false}
	}
	node.destinations = make([]string, 0)
	nodeStore[nodeName] = &node
	return &node
}

// processes "%lg -> zx, lx"
func parseLine(line string) (string, []string, string) {
	idx := strings.Index(line, "->")
	nodeType := line[0:1]
	nodeName := line[1 : idx-1]
	if nodeType == "b" {
		// keep the b of broadcaster
		nodeName = line[0 : idx-1]
	}
	destinations := strings.Split(line[idx+3:], ", ")
	basicNode := createNode(nodeName, nodeType)
	basicNode.destinations = append(basicNode.destinations, destinations...)

	return nodeName, destinations, nodeType
}

func parseInput(input []string) {
	var conNodes []string
	nodeMap := make(map[string][]string)
	for _, line := range input {
		node, destinations, nodeType := parseLine(line)
		// conjunction nodes need memory initialized
		if nodeType == "&" {
			conNodes = append(conNodes, node)
		}
		nodeMap[node] = destinations

	}
	rxNode := BasicNode{name: "rx", nodeType: Basic, isOn: false}
	nodeStore["rx"] = &rxNode

	// initialize memory for conjunction nodes, we need all incoming signals
	for _, conNodeName := range conNodes {
		conNode := nodeStore[conNodeName]
		var destinations []string
		for fromNodeName, _ := range nodeMap {
			// check if conNodeName is in nodeMap[fromNodeName]
			for _, dest := range nodeMap[fromNodeName] {
				if dest == conNodeName {
					destinations = append(destinations, fromNodeName)
				}
			}
		}
		conNode.initializeMemory(destinations)
	}
}

// -------------------------- Puzzle part 1 ----------------------------------

func SolvePart1() int {
	broadcaster := nodeStore["broadcaster"]
	for x := 0; x < 1000; x++ {
		buttonCount++
		broadcaster.receivePulse(broadcaster.name, Low)
		for len(messageQueue.messages) > 0 {
			msg := messageQueue.Pop()
			target := nodeStore[msg.to]
			target.receivePulse(msg.from.name, msg.pulse)
		}
	}
	lowPulseCount += buttonCount
	fmt.Printf("Low pulse count: %d, High pulse count: %d\n", lowPulseCount, highPulseCount)
	result := lowPulseCount * highPulseCount
	return result
}

// -------------------------- Puzzle part 2 ----------------------------------

var listOfNodes = []string{"tx", "kp", "gc", "vg"}
var countListOfNodes = make(map[string]int)

func Gcd(a, b int) int {
	for b != 0 {
		a, b = b, a%b
	}
	return a
}

func Lcm(a, b int) int {
	return a * b / Gcd(a, b)
}

// least common multiplier of countListOfNodes
func lcm() {
	var lcm int = 1
	for _, v := range countListOfNodes {
		lcm = Lcm(lcm, v)
	}
	fmt.Printf("Least common multiplier: %d\n", lcm)
}

func processStatus(target *BasicNode, pulse Pulse) {
	for _, node := range listOfNodes {
		if node == target.name {
			bqNode := nodeStore["bq"]
			if bqNode.memory[node] == High {
				countListOfNodes[node] = buttonCount
			}
		}
	}

	if len(countListOfNodes) == len(listOfNodes) {
		lcm()
		done = true
	}
}

func SolvePart2() int {
	broadcaster := nodeStore["broadcaster"]
	for x := 0; x < 1000000 && !done; x++ {
		buttonCount++
		broadcaster.receivePulse(broadcaster.name, Low)
		for len(messageQueue.messages) > 0 && !done {
			msg := messageQueue.Pop()
			target := nodeStore[msg.to]
			target.receivePulse(msg.from.name, msg.pulse)
			processStatus(target, msg.pulse)
		}
	}

	return 0
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
	result1 := SolvePart1()

	// reset state
	parseInput(input)
	buttonCount = 0
	lowPulseCount = 0
	highPulseCount = 0
	SolvePart2()
	stopwatch.Stop()

	// ---------------------- Print results ----------------------------------
	var elapsedTime time.Duration = stopwatch.GetElapsedTime()
	fmt.Println("Running as test:\t", runTest)
	fmt.Println("Result 1:\t\t", result1)
	fmt.Println("Elapsed time:\t\t", elapsedTime)
}
