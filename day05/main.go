package main

import (
	"fmt"
	"math"
	"strings"
	"time"

	"github.com/cdr74/AdventOfCode2023/utils"
)

// the flag runTest defines which data file to read
const runTest bool = false
const TEST_FILE string = "test.data"
const DATA_FILE string = "actual.data"

// ---------------------------------------------------------------------------

type Mapping struct {
	targetStart uint64
	sourceStart uint64
	sourceEnd   uint64
	length      uint64
}

type MappingList struct {
	mapping [][]Mapping
}

type Seed struct {
	start  uint64
	length uint64
}

func createMappings(input []string) MappingList {
	var mappings MappingList
	var currentMapping []Mapping

	for _, line := range input {

		if strings.HasSuffix(line, ":") {
			mappings.mapping = append(mappings.mapping, currentMapping)
			currentMapping = nil
			continue
		}

		if line == "" {
			continue
		}

		numbersList := strings.Fields(line)
		mapping := Mapping{
			targetStart: utils.StringToUint64(numbersList[0]),
			sourceStart: utils.StringToUint64(numbersList[1]),
			length:      utils.StringToUint64(numbersList[2]),
			sourceEnd:   utils.StringToUint64(numbersList[1]) + utils.StringToUint64(numbersList[2]),
		}

		currentMapping = append(currentMapping, mapping)
	}
	mappings.mapping = append(mappings.mapping, currentMapping)
	return mappings
}

func applyMapping(position uint64, mappings []Mapping) uint64 {
	for _, mapping := range mappings {
		if position >= mapping.sourceStart && position <= mapping.sourceEnd {
			delta := position - mapping.sourceStart
			//fmt.Printf("applyMapping() - %v -> %v -> %v\n", position, mappings, mapping.targetStart+delta)
			return mapping.targetStart + delta
		}
	}
	// no mapping keep as is
	//fmt.Printf("applyMapping() - %v -> %v -> %v\n", position, mappings, position)
	return position
}

// ---------------------------------------------------------------------------

func getSeeds(input string) []uint64 {
	var seeds []uint64
	numbersList := strings.Split(input, " ")
	for _, number := range numbersList {
		if number != "" {
			seeds = append(seeds, utils.StringToUint64(strings.Trim(number, " ")))
		}
	}
	fmt.Printf("getSeeds() - %v\n", seeds)
	return seeds
}

func SolvePuzzle1(seeds []uint64, mappings MappingList) uint64 {
	var result uint64 = math.MaxUint64
	var position uint64 = 0

	for _, seed := range seeds {
		position = seed
		for mappingLevel := 0; mappingLevel < 7; mappingLevel++ {
			position = applyMapping(position, mappings.mapping[mappingLevel])
		}
		if position < result {
			result = position
		}
		//fmt.Printf("SolvePuzzle1() - Seed %v -> Position %v\n", seed, position)
	}

	return result
}

// ---------------------------------------------------------------------------

func getSeeds2(input string) []Seed {
	var seedsList []Seed
	numbersList := strings.Fields(input)

	for i := 0; i < len(numbersList); i += 2 {
		start := utils.StringToUint64(numbersList[i])
		length := utils.StringToUint64(numbersList[i+1])
		seedsList = append(seedsList, Seed{start: start, length: length})
	}
	fmt.Printf("getSeeds2() - %v\n", seedsList)
	return seedsList
}

func SolvePuzzle2(seeds []Seed, mappings MappingList) uint64 {
	var result uint64 = math.MaxUint64
	var position uint64 = 0

	for _, seedRange := range seeds {
		fmt.Printf("> Working on: %v\n", seedRange)
		for seed := seedRange.start; seed < seedRange.start+seedRange.length; seed++ {
			position = seed
			for mappingLevel := 0; mappingLevel < 7; mappingLevel++ {
				position = applyMapping(position, mappings.mapping[mappingLevel])
			}
			if position < result {
				result = position
			}
			//fmt.Printf("SolvePuzzle2() - Seed %v -> Position %v\n", seed, position)
		}
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

	var seeds []uint64
	var seedsString string
	idx := strings.IndexRune(input[0], ':')
	seedsString = input[0][idx+1:]
	seeds = getSeeds(seedsString)

	// consume first 3 lines
	input = input[3:]
	mappings := createMappings(input)

	result1 := SolvePuzzle1(seeds, mappings)

	seeds2 := getSeeds2(seedsString)
	result2 := SolvePuzzle2(seeds2, mappings)
	stopwatch.Stop()

	// -------------------------------------
	var elapsedTime time.Duration = stopwatch.GetElapsedTime()
	fmt.Println("Running as test:\t", runTest)
	fmt.Println("Result 1:\t\t\t", result1)
	fmt.Println("Result 2:\t\t\t", result2)
	fmt.Println("Elapsed time:\t\t", elapsedTime)
}
