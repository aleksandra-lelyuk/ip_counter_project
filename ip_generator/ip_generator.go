package ip_generator

import (
	"archive/zip"
	"fmt"
	"io"
	"math/rand"
	"os"
	"time"
)

func generateRandomIPv4() string {
	// Generate a random IPv4 address
	return fmt.Sprintf("%d.%d.%d.%d", rand.Intn(256), rand.Intn(256), rand.Intn(256), rand.Intn(256))
}

func GenerateIPFile(filename string, N int, K int) error {
	// Ensure K is not greater than N
	if K > N {
		K = N
	}

	// Use a slice to store unique addresses
	uniqueIPs := make([]string, 0, K)

	// Seed the random number generator
	rand.Seed(time.Now().UnixNano())

	// Generate K unique addresses
	uniqueMap := make(map[string]struct{})
	for len(uniqueIPs) < K {
		ip := generateRandomIPv4()
		if _, exists := uniqueMap[ip]; !exists {
			uniqueIPs = append(uniqueIPs, ip)
			uniqueMap[ip] = struct{}{}
		}
	}

	// Shuffle unique IPs
	rand.Shuffle(len(uniqueIPs), func(i, j int) {
		uniqueIPs[i], uniqueIPs[j] = uniqueIPs[j], uniqueIPs[i]
	})

	// Create or open the output file
	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("error creating file: %w", err)
	}
	defer file.Close()

	// Write the unique IPs to the file
	for _, ip := range uniqueIPs {
		file.WriteString(ip + "\n")
	}

	// Fill the rest of the addresses with copies of the unique ones
	for i := K; i < N; i++ {
		file.WriteString(uniqueIPs[rand.Intn(K)] + "\n")
	}

	fmt.Printf("Successfully generated %d random IPv4 addresses with %d unique in %s\n",
		N, K, filename)
	return nil
}

func ZipFile(sourceFile, zipFileName string) error {
	// Create the ZIP file
	zipFile, err := os.Create(zipFileName)
	if err != nil {
		return fmt.Errorf("failed to create zip file: %v", err)
	}
	// Use defer to ensure that zipFile.Close() is always called,
	defer func() {
		if err := zipFile.Close(); err != nil {
			fmt.Println("Error closing zip file:", err)
		}
	}()

	// Create a new zip writer
	zipWriter := zip.NewWriter(zipFile)
	defer func() {
		if err := zipWriter.Close(); err != nil {
			fmt.Println("Error closing zip writer:", err)
		}
	}()

	// Open the source file to be zipped
	srcFile, err := os.Open(sourceFile)
	if err != nil {
		return fmt.Errorf("failed to open source file: %v", err)
	}
	defer func() {
		if err := srcFile.Close(); err != nil {
			fmt.Println("Error closing source file:", err)
		}
	}()

	// Create a writer inside the zip for the file
	zipEntry, err := zipWriter.Create(sourceFile)
	if err != nil {
		return fmt.Errorf("failed to create zip entry: %v", err)
	}

	// Copy the content of the source file into the zip entry
	if _, err = io.Copy(zipEntry, srcFile); err != nil {
		return fmt.Errorf("failed to write file to zip: %v", err)
	}

	return nil
}
