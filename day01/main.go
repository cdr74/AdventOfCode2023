package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"time"
)

// the flag runTest defines which data file to read
const runTest bool = true
const TEST_FILE string = "test.data"
const DATA_FILE string = "actual.data"

type Stopwatch struct {
	startTime   time.Time
	elapsedTime time.Duration
	isRunning   bool
}

func NewStopwatch() *Stopwatch {
	return &Stopwatch{
		startTime:   time.Now(),
		elapsedTime: 0,
		isRunning:   false,
	}
}

func (s *Stopwatch) Start() {
	s.isRunning = true
	s.startTime = time.Now()
}

func (s *Stopwatch) Stop() {
	s.isRunning = false
	s.elapsedTime = time.Since(s.startTime)
}

func readDataFile(filename string) []string {
	file, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	// Create a scanner.
	scanner := bufio.NewScanner(file)

	// Create an array to store the lines.
	var lines []string

	// Scan the file line by line.
	for scanner.Scan() {
		line := scanner.Text()
		lines = append(lines, line)
	}

	// Return the array of lines.
	return lines
}

func solvePuzzle() string {
	return "tbd"
}

func main() {
	var result string
	stopwatch := NewStopwatch()
	stopwatch.Start()

	if runTest {
		readDataFile(TEST_FILE)
	} else {
		readDataFile(DATA_FILE)
	}

	result = solvePuzzle()
	stopwatch.Stop()

	// -------------------------------------
	elapsedTime := stopwatch.elapsedTime
	fmt.Println("Running as test:\t", runTest)
	fmt.Println("Result:\t\t\t", result)
	fmt.Println("Elapsed time:\t\t", elapsedTime)
}
