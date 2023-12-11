package utils

func CreateIntArray(width int, depth int, initial int) [][]int {
	dataArr := make([][]int, depth)

	for i := range dataArr {
		dataArr[i] = make([]int, width)
		for j := range dataArr[i] {
			dataArr[i][j] = initial
		}
	}
	return dataArr
}

func SumOfArray(array []int) int {
	sum := 0
	for _, num := range array {
		sum += num
	}
	return sum
}

func Min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func Max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func Abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}
