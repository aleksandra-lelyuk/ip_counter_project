package main

import (
	"fmt"
	"ip_counter_project/pkg/ip_generator"
)

func main() {
	fileNames := []string{"small.txt", "medium.txt", "large.txt"}
	numEntries := []int{10_000_000, 100_000_000, 300_000_000}
	numUnique := []int{2_000_000, 20_000_000, 60_000_000} //note potential bias because num unique is always 20% of total entries

	for i := 0; i < 3; i++ {
		err := ip_generator.GenerateIPFile(fileNames[i], numEntries[i], numUnique[i])
		if err != nil {
			fmt.Println("Error:", err)
		}
	}
}
