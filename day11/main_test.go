package main

import "testing"

func TestCountOfValuesBetween(t *testing.T) {
	values := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}

	// Test case 1: start = 2, end = 8
	expected := 5
	result := countOfValuesBetween(2, 8, values)
	if result != expected {
		t.Errorf("countOfValuesBetween(2, 8, values) = %d; want %d", result, expected)
	}

	// Test case 2: start = 5, end = 10
	expected = 4
	result = countOfValuesBetween(5, 10, values)
	if result != expected {
		t.Errorf("countOfValuesBetween(5, 10, values) = %d; want %d", result, expected)
	}

	// Test case 3: start = 1, end = 1
	expected = 0
	result = countOfValuesBetween(1, 1, values)
	if result != expected {
		t.Errorf("countOfValuesBetween(1, 1, values) = %d; want %d", result, expected)
	}
}
