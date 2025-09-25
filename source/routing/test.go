package main

import (
	"fmt"
)

func printSlice(s []int) {
	fmt.Printf("len=%d cap=%d %v\n", len(s), cap(s), s)
}

func _sublineInPolygon(quantizedLine []int, target int, start, end int) bool {
	// Base case
	if start > end {
		return false
	}
	
	middle := (start + end) / 2
	fmt.Println("MIDDLE")
	fmt.Printf("start=%d end=%d middle=%d value=%d\n\n", start, end, middle, quantizedLine[middle]) 

	if quantizedLine[middle] == target {
		return true
	}

	// Here we need to left and right search
	fmt.Println("LEFT")
	if _sublineInPolygon(quantizedLine, target, start, middle-1) {
		return true
	}

	fmt.Println("RIGHT")
	if _sublineInPolygon(quantizedLine, target, middle+1, end) {
		return true
	}

	return false
}

func main() {
	s := []int{2, 3, 5, 7, 11, 13, 20, 15, 10, 6, 55}
	printSlice(s)

	start := 0
	end := len(s)-1

	_sublineInPolygon(s, 14, start, end)
}
