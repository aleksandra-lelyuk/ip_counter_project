package main

import (
	"fmt"
	"ip_counter_project/experiment_setup"
	"ip_counter_project/print_results"
	"time"
)

func main() {

	start := time.Now() // Start timing

	//choose which files to use
	fileName := "ip_addresses"
	parseIP := false                //choose whether to parse IPs before analysis or not
	hashTypes := []string{"xxhash"} //select available hash functions ["fnv32", "fnv64", "xxhash"]
	numHashesList := []int{16}      //choose number of hash functions (registers) to use

	//Initialize csv file for results
	outputFileName := "../../data/final_result_xxhash.csv"
	header := []string{
		"fileSize", "numHashes", "hashType", "parseIP",
		"averageEstimatedCount", "harmonicEstimatedCount", "customEstimatedCount",
		"averageOfAverages", "elapsedTimeSec", "realCount",
	}

	experiment_setup.InitializeCsvWithHeader(outputFileName, header)

	for _, hashType := range hashTypes {
		for _, numHashes := range numHashesList {
			experiment_setup.RunExperiment(fileName, outputFileName, hashType, numHashes, parseIP, 1_000_000_000)
		}
	}
	elapsed := time.Since(start) // End timing
	minutes := int(elapsed.Minutes())
	seconds := int(elapsed.Seconds()) % 60
	fmt.Printf("Total time taken: %d minutes and %d seconds\n", minutes, seconds)

	print_results.PrintResultDataExperiment(outputFileName)
}
