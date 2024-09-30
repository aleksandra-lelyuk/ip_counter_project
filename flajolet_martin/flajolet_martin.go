package flajolet_martin

import (
	"bufio"
	"fmt"
	"github.com/cespare/xxhash/v2"
	"hash/fnv"
	"log"
	"math"
	"net"
	"os"
	"strings"
	"sync"
)

// Function to count the number of trailing zeros in the binary representation of an integer
func CountTrailingZeros(n interface{}) int {
	switch v := n.(type) {
	case uint32:
		if v == 0 {
			return 0 //avoid extreme overestimation
		}
		count := 0
		for (v & 1) == 0 { // Check if the righttmost digit == 0
			count++
			v >>= 1
		}
		return count
	case uint64:
		if v == 0 {
			return 0 //avoid extreme overestimation
		}
		count := 0
		for (v & 1) == 0 {
			count++
			v >>= 1
		}
		return count
	default:
		log.Fatalf("Unsupported type: %T", n)
		return -1
	}
}

// Define a type for the hash function to allow dynamic hash function selection
type HashFunction func([]byte, byte) uint64

// Function to select the hash function based on user input
func selectHashFunction(hashType string) HashFunction {
	switch hashType {
	case "fnv64":
		return func(value []byte, seed byte) uint64 {
			h := fnv.New64a()
			h.Write(value)
			h.Write([]byte{seed})
			return h.Sum64()
		}
	case "fnv32":
		return func(value []byte, seed byte) uint64 {
			h := fnv.New32a()
			h.Write(value)
			h.Write([]byte{seed})
			return uint64(h.Sum32())
		}
	case "xxhash":
		return func(value []byte, seed byte) uint64 {
			h := xxhash.New()
			h.Write(value)
			h.Write([]byte{seed})
			return h.Sum64()
		}
	default:
		log.Fatalf("Unsupported hash type: %s", hashType)
		return nil
	}
}

func CreateParallelFJByteStreams(txtFileName string, numHashes int, hashType string, parseIP bool) (float64, float64, float64, float64, error) {
	//Print stats for this test case
	fmt.Printf("\n\n%s\nCONFIGURATIONS\n%s\n", strings.Repeat("-", 50), strings.Repeat("-", 50))
	fmt.Printf("File Name: %s\n", txtFileName)
	fmt.Printf("Number of Hashes: %d\n", numHashes)
	fmt.Printf("Hash Type: %s\n", hashType)
	fmt.Printf("Parse IPs: %t\n", parseIP)
	fmt.Printf("%s\n", strings.Repeat("-", 50))

	trailingZeros := make([]int, numHashes)

	var wg sync.WaitGroup
	resultChan := make(chan struct {
		index int
		value int
	}, numHashes)

	hashFunc := selectHashFunction(hashType) // Set hash function

	// Parallelize complete algorithm instead of just hash calculations to take most advantage of parallelization
	for i := 0; i < numHashes; i++ { // Run FJ algorithm in parallel for each hash function (including reading from file)
		wg.Add(1)
		go func(index int) {
			defer wg.Done()
			value, err := FlajoletMartinStreamBytesSingle(txtFileName, index, hashFunc, parseIP)
			if err != nil {
				fmt.Printf("Error: %v\n", err)
				return
			}
			resultChan <- struct {
				index int
				value int
			}{index, value}
		}(i)
	}

	go func() {
		wg.Wait()
		close(resultChan)
	}()

	for result := range resultChan {
		trailingZeros[result.index] = result.value
	}

	fmt.Println("Number of trailing zeros from each hash:")
	fmt.Println(trailingZeros)

	averageEstimatedCount := CalculateRegularMean(trailingZeros, numHashes)
	fmt.Printf("Average Estimated Count: %.0f\n", averageEstimatedCount)

	harmonicEstimatedCount := CalculateHarmonicMean(trailingZeros, numHashes)
	fmt.Printf("Harmonic Estimated Count: %.0f\n", harmonicEstimatedCount)

	customEstimatedCount := CalculateCustomMean(trailingZeros, numHashes)
	fmt.Printf("Custom Estimated Count: %.0f\n", customEstimatedCount)

	averageOfAverages := (averageEstimatedCount + harmonicEstimatedCount + customEstimatedCount) / 3

	return averageEstimatedCount, harmonicEstimatedCount, customEstimatedCount, averageOfAverages, nil
}

func FlajoletMartinStreamBytesSingle(txtFilePath string, thisHashNum int, hashFunc HashFunction, parseIP bool) (int, error) {
	//seed := byte(rand.Intn(100)) // Random seed between 0 and 99
	trailingZeroCount := 0
	lineCount := 0

	file, err := os.Open(txtFilePath)
	if err != nil {
		return 0, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	//fmt.Println("Reading with buffer")
	buf := make([]byte, 256*1024) // 256 KB buffer to speed upt reading
	scanner.Buffer(buf, 256*1024)

	for scanner.Scan() {
		lineCount++ // Increment line counter
		if thisHashNum == 0 && lineCount%1_000_000_000 == 0 {
			fmt.Printf("Hash 0: Processed %d billion lines\n", lineCount/1_000_000_000)
		}

		var elemBytes []byte

		if parseIP == false {
			elemBytes = scanner.Bytes()
		} else {
			elemString := scanner.Text()
			ip := net.ParseIP(elemString)
			if ip == nil {
				fmt.Println("Invalid IP address")
				continue
			}

			elemBytes = ip.To4()

		}
		// Make hashing function pseudorandom by writing the ID of the hash function as seed
		hashedValue := hashFunc(elemBytes, byte(thisHashNum))                       // Get the hash value using the selected hash function
		trailingZeroCount = max(trailingZeroCount, CountTrailingZeros(hashedValue)) //Count trailing zeros in the binary representation of hash value

	}
	if thisHashNum == 0 {
		fmt.Printf("Processed %d lines\n", lineCount)
	}

	return trailingZeroCount, nil
}

// Function to get the maximum of two integers
func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func CalculateHarmonicMean(trailingZeros []int, numHashes int) float64 {
	sum := 0.0
	for _, zeros := range trailingZeros {
		sum += 1 / math.Pow(2, float64(zeros)) //harmonic average
	}
	average := float64(numHashes) / sum //harmonic average
	alpha := 0.77351

	estimatedCount := alpha * average
	return estimatedCount
}

func CalculateRegularMean(trailingZeros []int, numHashes int) float64 {
	sum := 0.0
	for _, zeros := range trailingZeros {
		sum += math.Pow(2, float64(zeros)) //regular average
	}
	average := sum / float64(numHashes) //regular average
	alpha := 0.77351
	//alpha := 1.0
	//correctedEstimate := alpha * average
	estimatedCount := alpha * average
	return estimatedCount
}

func CalculateCustomMean(trailingZeros []int, numHashes int) float64 {
	sum := 0
	for _, zeros := range trailingZeros {
		sum += zeros
	}
	average := float64(sum) / float64(numHashes)
	alpha := 0.77351

	estimatedCount := alpha * math.Pow(2, average)
	return estimatedCount
}
