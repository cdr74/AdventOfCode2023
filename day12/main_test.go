package main

import "testing"

func TestMemoizedRecursiveCount(t *testing.T) {
	line := []byte("..#.")
	sequenceList := []int{1}

	expected := 1
	result := memoizedRecursiveCount(line, sequenceList)
	if result != expected {
		t.Errorf("Expected %d, but got %d", expected, result)
	}

	line = []byte("##..#.")
	sequenceList = []int{2, 1}

	expected = 1
	result = memoizedRecursiveCount(line, sequenceList)
	if result != expected {
		t.Errorf("Expected %d, but got %d", expected, result)
	}

	line = []byte("??..#.")
	sequenceList = []int{2, 1}

	expected = 1
	result = memoizedRecursiveCount(line, sequenceList)
	if result != expected {
		t.Errorf("Expected %d, but got %d", expected, result)
	}

	line = []byte(".???.#.")
	sequenceList = []int{2, 1}

	expected = 2
	result = memoizedRecursiveCount(line, sequenceList)
	if result != expected {
		t.Errorf("Expected %d, but got %d", expected, result)
	}

	line = []byte("?.##.")
	sequenceList = []int{1, 2}

	expected = 1
	result = memoizedRecursiveCount(line, sequenceList)
	if result != expected {
		t.Errorf("Expected %d, but got %d", expected, result)
	}
}
