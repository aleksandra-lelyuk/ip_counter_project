package bitstore

import (
	"bufio"
	"encoding/binary"
	"fmt"
	"net"
	"os"
)

func CountUniqueIPs(txtFilePath string) (int, error) {
	file, err := os.Open(txtFilePath)
	if err != nil {
		return 0, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	buf := make([]byte, 256*1024) // 256 KB buffer
	scanner.Buffer(buf, 256*1024)

	// 2^32 bits can be used to represent all possible IPv4 addresses (each bit mapped to one address)
	// To initialize the slice of correct size, shift the binary representation of 1 by 32 positions
	// Divide by 8 to determine the number of bytes needed to store
	bytesNeeded := 1 << 32 / 8 //=2^{29} bytes = 512 MB.
	bitset := make([]byte, bytesNeeded)

	for scanner.Scan() {

		elemString := scanner.Text()
		ip := net.ParseIP(elemString)

		if ip == nil {
			fmt.Println("Invalid IP address")
			continue
		}

		ipBytes := ip.To4()

		ipUint := binary.BigEndian.Uint32(ipBytes) //convert bytes to unsigned integer
		// Find the corresponding byte (bitstore[ipUint/8]) and then the bit within that byte (ipUint % 8)
		byteIndex := ipUint / 8
		bitIndex := ipUint % 8
		// Since 1 << bitIndex creates an 8 bit representation with only 1 non-zero bit, the |= can be used
		bitset[byteIndex] |= 1 << bitIndex // |= operator sets corresponding bit to 1 within the byte

	}
	uniqueCount := 0
	for _, byteVal := range bitset { // go through bitstring byte by byte
		uniqueCount += countSetBits(byteVal)
	}

	return uniqueCount, nil
}

func countSetBits(b byte) int {
	count := 0
	for b > 0 { //while value of byte>0
		count += int(b & 1) //add 1 to count if least significant bit if b == 1
		b >>= 1             //discard least significant bit in b
	}
	return count
}
