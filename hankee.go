package main

import "fmt"

func main() {
	cpu := NewCPU()
	cpu.loadAndRun([]uint8{1})
	fmt.Println("Hello, World!")
}
