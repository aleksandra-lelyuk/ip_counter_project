package main

import (
	"fmt"
	"ip_counter_project/pkg/experiment_setup"
	"ip_counter_project/pkg/print_results"
	"time"
)

func main() {

	start := time.Now() // Start timing

	fileMap := map[string]string{
		"small":  "small.txt",
		"medium": "medium.txt",
		"large":  "large.txt",
	}

	realCountMap := map[string]int{
		"small":  2_000_000,
		"medium": 20_000_000,
		"large":  60_000_000,
	}

	fileSizes := []string{"small"}                    //choose which files to use
	parseIPs := []bool{false}                         //choose whether to parse IPs before analysis or not
	hashTypes := []string{"fnv32", "fnv64", "xxhash"} //select available hash functions ["fnv32", "fnv64", "xxhash"]
	numHashesList := []int{8, 16, 64, 128}            //choose number of hash functions (registers) to use

	//Initialize csv file for results
	outputFileName := "../../data/print_test_result.csv"
	header := []string{
		"fileSize", "numHashes", "hashType", "parseIP",
		"averageEstimatedCount", "harmonicEstimatedCount", "customEstimatedCount",
		"averageOfAverages", "elapsedTimeSec", "realCount",
	}

	experiment_setup.InitializeCsvWithHeader(outputFileName, header)

	for _, hashType := range hashTypes {
		for _, numHashes := range numHashesList {
			for _, fileSize := range fileSizes {
				for _, parseIP := range parseIPs {
					experiment_setup.RunExperiment(fileMap[fileSize], outputFileName, hashType, numHashes, parseIP, realCountMap[fileSize])
				}
			}
		}
	}
	elapsed := time.Since(start) // End timing
	minutes := int(elapsed.Minutes())
	seconds := int(elapsed.Seconds()) % 60
	fmt.Printf("Total time taken: %d minutes and %d seconds\n", minutes, seconds)

	print_results.PrintResultDataExperiment(outputFileName)
}
