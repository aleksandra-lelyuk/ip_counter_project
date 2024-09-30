package main

import (
	"fmt"
	"ip_counter_project/bitstore"
	"os"
	"runtime/pprof"
	"time"
)

func main() {
	testDatasets := []string{"small", "medium", "large"}

	for _, dataset := range testDatasets {
		start := time.Now()
		txtFilePath := fmt.Sprintf("../flajolet_martin_experiment/%s.txt", dataset)
		fmt.Printf("Processing file: %s\n", txtFilePath)
		numUnique, _ := bitstore.CountUniqueIPs(txtFilePath)
		fmt.Printf("Number of unique elements: %d\n", numUnique)
		elapsed := time.Since(start) // End timing
		minutes := int(elapsed.Minutes())
		seconds := int(elapsed.Seconds()) % 60
		fmt.Printf("Total time taken by bitstore algo on %s test dataset: %d minutes and %d seconds\n", dataset, minutes, seconds)
	}

	profileName := fmt.Sprintf("../../data/cpu_profiles/cpu_profile_bitstore_final.prof")
	f, err := os.Create(profileName)
	if err != nil {
		fmt.Println("could not create CPU profile: ", err)
		return
	}
	defer f.Close()

	// Start CPU profiling
	if err := pprof.StartCPUProfile(f); err != nil {
		fmt.Println("could not start CPU profile: ", err)
		return
	}
	defer pprof.StopCPUProfile()

	start := time.Now()
	txtFilePath := "../flajolet_martin/ip_addresses"
	fmt.Printf("Processing file: %s\n", txtFilePath)
	numUnique, _ := bitstore.CountUniqueIPs(txtFilePath)
	fmt.Printf("Number of unique elements: %d\n", numUnique)
	elapsed := time.Since(start) // End timing
	minutes := int(elapsed.Minutes())
	seconds := int(elapsed.Seconds()) % 60
	fmt.Printf("Total time taken by bitstore algo on FINAL dataset: %d minutes and %d seconds\n", minutes, seconds)
}
