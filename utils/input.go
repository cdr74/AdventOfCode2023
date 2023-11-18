package utils

import (
	"bufio"
	"log"
	"os"
)

func ReadDataFile(filename string) []string {
	file, err := os.Open(filename)
	if err != nil {
		// log.Fates does exit the prog
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	
	var lines []string
	for scanner.Scan() {
		line := scanner.Text()
		lines = append(lines, line)
	}

	return lines
}

