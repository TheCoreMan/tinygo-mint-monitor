package main

import (
	"fmt"
	"machine"
	"time"
)

func main() {
	count := 0
	freq := machine.CPUFrequency()
	for {
		fmt.Printf("%d: It's a-me!\n", count)
		n, err := machine.GetRNG()
		if err != nil {
			fmt.Printf("Error getting a random number: %v\n", err)
		}
		fmt.Printf("Random number: %d\n", n)
		fmt.Printf("CPU Frequency: %d Hz\n", freq)
		fmt.Printf("The time now is: %v\n", time.Now())
		time.Sleep(time.Second)
		count++
	}
}
