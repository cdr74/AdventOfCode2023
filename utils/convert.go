package utils

import (
	"log"
	"strconv"
)

func StringToInt(s string) int {
  i, err := strconv.ParseInt(s, 10, 0)  
  if err != nil {
    // log.Fates does exit the prog
    log.Fatal(err)
  }
  return int(i)
}

func IntToString(i int) string {
	s := strconv.Itoa(i)
	return s
}

// assumes each digit is an int
func StringToIntArray(input string) []int {
	intArray := make([]int, 0, len(input))

	for _, char := range input {
    num := StringToInt(string(char))
		intArray = append(intArray, num)
	}

	return intArray
}

func StringToByteArray(s string) []byte {  
  return []byte(s)
}