package experiment_setup

import (
	"encoding/csv"
	"fmt"
	"ip_counter_project/flajolet_martin"
	"os"
	"runtime/pprof"
	"strconv"
	"strings"
	"time"
)

func RunExperiment(txtFileName string, outputFileName string, hashType string, numHashes int, parseIP bool, realCount int) []string {
	var parseIPStr string
	if parseIP {
		parseIPStr = "T"
	} else {
		parseIPStr = "F"
	}

	profileName := fmt.Sprintf("../../data/cpu_profiles/cpu_profile_%s_%s_%d_%s.prof", hashType, parseIPStr, numHashes, txtFileName[:len(txtFileName)-4])
	f, err := os.Create(profileName)
	if err != nil {
		fmt.Println("could not create CPU profile: ", err)
		return nil
	}
	defer f.Close()

	// Start CPU profiling
	if err := pprof.StartCPUProfile(f); err != nil {
		fmt.Println("could not start CPU profile: ", err)
		return nil
	}
	defer pprof.StopCPUProfile()

	start := time.Now()

	averageEstimatedCount, harmonicEstimatedCount, customEstimatedCount, averageOfAverages, _ :=
		flajolet_martin.CreateParallelFJByteStreams(txtFileName, numHashes, hashType, parseIP)
	fmt.Printf("Average estimation: %.0f\n", averageOfAverages)
	elapsed := time.Since(start)
	elapsedSeconds := elapsed.Seconds()

	fmt.Printf("Time taken: %.2f s\n", elapsedSeconds)

	// Save input and output parameters to a file
	statsValues := []string{
		outputFileName,
		txtFileName,
		strconv.Itoa(numHashes),
		hashType,
		strconv.FormatBool(parseIP),
		strconv.FormatFloat(averageEstimatedCount, 'f', 0, 64),
		strconv.FormatFloat(harmonicEstimatedCount, 'f', 0, 64),
		strconv.FormatFloat(customEstimatedCount, 'f', 0, 64),
		strconv.FormatFloat(averageOfAverages, 'f', 0, 64),
		strconv.FormatFloat(elapsedSeconds, 'f', 2, 64),
		strconv.Itoa(realCount),
	}

	fmt.Println(statsValues)
	fmt.Printf("%s\n", strings.Repeat("-", 50))

	if err := writeStatsToCsv(outputFileName, statsValues); err != nil {
		fmt.Printf("Error writing stats to CSV file: %v\n", err)
		return nil
	}
	return statsValues
}

func writeStatsToCsv(outputFileName string, statsValues []string) error {
	file, err := os.OpenFile(outputFileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("could not open file: %v", err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	if err := writer.Write(statsValues); err != nil {
		return fmt.Errorf("could not write stats values to CSV file: %v", err)
	}

	return nil
}

func InitializeCsvWithHeader(outputFileName string, header []string) {

	file, err := os.OpenFile(outputFileName, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		fmt.Printf("Error opening file %s: %v\n", outputFileName, err)
		return
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	if err := writer.Write(header); err != nil {
		fmt.Printf("Error writing header to CSV file %s: %v\n", outputFileName, err)
		return
	}
}
