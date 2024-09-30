package main

import (
	"ip_counter_project/print_results"
)

func main() {
	resultsCsv := "../../data/large_test_result.csv"
	outputFile := "../../data/graphs/large_test_res.png"
	print_results.PlotResultDataExperiment(resultsCsv, outputFile)
}
