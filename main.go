// Package main implements a command-line utility for expanding IP addresses and CIDR ranges.
package main

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"os"

	"github.com/spf13/pflag"
)

func main() {
	// Define command-line flags
	printIPv4 := pflag.BoolP("ipv4", "4", false, "Print only IPv4 addresses")
	printIPv6 := pflag.BoolP("ipv6", "6", false, "Print only IPv6 addresses")
	help := pflag.BoolP("help", "h", false, "Display this help message")

	pflag.Parse()

	// If help flag is set, print usage
	if *help {
		printUsage()
		return
	}

	// Determine IP address filtering based on flags
	// If neither or both flags are set, include both IPv4 and IPv6
	includeIPv4 := *printIPv4 || !(*printIPv4) && !(*printIPv6)
	includeIPv6 := *printIPv6 || !(*printIPv4) && !(*printIPv6)

	// Determine input source: file if provided, otherwise stdin
	var reader io.Reader
	if len(pflag.Args()) > 0 {
		filename := pflag.Args()[0]
		file, err := os.Open(filename)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		defer file.Close()
		reader = file
	} else {
		reader = os.Stdin
	}

	processInput(reader, includeIPv4, includeIPv6)
}

// processInput reads from the provided reader and processes each line
// to extract and print IP addresses based on the specified filters.
func processInput(reader io.Reader, includeIPv4, includeIPv6 bool) {
	scanner := bufio.NewScanner(reader)
	for scanner.Scan() {
		line := scanner.Text()
		printIPsFromLine(line, includeIPv4, includeIPv6)
	}

	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

// printIPsFromLine parses a single line as an IP address or CIDR range
// and prints the contained IP addresses based on the specified filters.
func printIPsFromLine(line string, includeIPv4, includeIPv6 bool) {
	// First, try parsing as a single IP address
	ip := net.ParseIP(line)
	if ip != nil {
		if (includeIPv4 && ip.To4() != nil) || (includeIPv6 && ip.To4() == nil) {
			fmt.Println(ip)
		}
		return
	}

	// If not a single IP, try parsing as a CIDR range
	_, ipNet, err := net.ParseCIDR(line)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return
	}

	// Iterate through all IPs in the CIDR range
	for ip := ipNet.IP.Mask(ipNet.Mask); ipNet.Contains(ip); incrementIP(ip) {
		if (includeIPv4 && ip.To4() != nil) || (includeIPv6 && ip.To4() == nil) {
			fmt.Println(ip)
		}
	}
}

// incrementIP increments the given IP address by 1.
// It properly handles overflow across octets/hexadecets.
func incrementIP(ip net.IP) {
	// Traverse the IP address from the least significant byte (rightmost) to the
	// most significant byte (leftmost).
	for j := len(ip) - 1; j >= 0; j-- {
		ip[j]++ // Increment the current byte by 1.

		// If the incremented byte is not zero, there was no carry-over, and we
		// can exit the loop. If it is zero, it means the increment caused an
		// overflow (e.g., from 0xFF to 0x00), and we need to increment the next
		// more significant byte.
		if ip[j] != 0 {
			break
		}
	}
}

// printUsage displays the program usage information.
func printUsage() {
	fmt.Println("cidrex - Expand CIDR ranges")
	fmt.Println("\nUsage:")
	fmt.Println("  cidrex [OPTIONS] [filename]")
	fmt.Println("\nOptions:")
	pflag.PrintDefaults()
	fmt.Println("\nExamples:")
	fmt.Println("  cidrex input.txt")
	fmt.Println("  cidrex -4 input.txt")
	fmt.Println("  cat input.txt | cidrex -6")
}
