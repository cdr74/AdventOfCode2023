package utils

func SumOfArray(array []int) int {
	sum := 0
  	for _, num := range array {
    	sum += num
  	}
	return sum
}

