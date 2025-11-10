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

type Int interface {
	Func()
	InnerFunc()
}

type ClassA struct {}

func (c *ClassA) Func() {
	fmt.Printf("ClassA - Func\n")
	c.InnerFunc()
}

func (c *ClassA) InnerFunc() {
	fmt.Printf("ClassA - InnerFunc\n")
}

type ClassB struct {
	*ClassA
}

func (c *ClassB) InnerFunc() {
	fmt.Printf("ClassB - InnerFunc\n")
}

func main() {
	a := &ClassA{}
	b := &ClassB{
		ClassA: &ClassA{},
	}

	var intA Int = a
	var intB Int = b

	intA.Func()
	intB.Func()
}
