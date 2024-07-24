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

	// Create a new buffered writer to stdout
	var writer = bufio.NewWriterSize(os.Stdout, 32*1024)
	defer writer.Flush()

	if err := processInput(reader, writer, includeIPv4, includeIPv6); err != nil {
		fmt.Fprintf(os.Stderr, "Error processing input: %v\n", err)
		os.Exit(1)
	}
}

// processInput reads from the provided reader and processes each line
// to extract and print IP addresses based on the specified filters.
func processInput(reader io.Reader, writer io.Writer, includeIPv4, includeIPv6 bool) error {
	scanner := bufio.NewScanner(reader)
	for scanner.Scan() {
		if err := printIPsFromLine(writer, scanner.Text(), includeIPv4, includeIPv6); err != nil {
			return err
		}
	}

	return scanner.Err()
}

// printIPsFromLine parses a single line as an IP address or CIDR range
// and prints the contained IP addresses based on the specified filters.
func printIPsFromLine(writer io.Writer, line string, includeIPv4, includeIPv6 bool) error {
	// First, try parsing as a single IP address
	ip := net.ParseIP(line)
	if ip != nil {
		return printIP(writer, ip, includeIPv4, includeIPv6)
	}

	// If not a single IP, try parsing as a CIDR range
	_, ipNet, err := net.ParseCIDR(line)
	if err != nil {
		// Print message to stderr but don't return an error to continue processing
		fmt.Fprintf(os.Stderr, "invalid IP or CIDR: %s\n", line)
		return nil
	}

	// Iterate through all IPs in the CIDR range
	for ip := ipNet.IP.Mask(ipNet.Mask); ipNet.Contains(ip); incrementIP(ip) {
		if err := printIP(writer, ip, includeIPv4, includeIPv6); err != nil {
			return err
		}
	}

	return nil
}

// printIP writes the given IP address to the provided writer if it matches the
// inclusion criteria specified by includeIPv4 and includeIPv6.
func printIP(writer io.Writer, ip net.IP, includeIPv4, includeIPv6 bool) error {
	if (includeIPv4 && ip.To4() != nil) || (includeIPv6 && ip.To4() == nil) {
		_, err := fmt.Fprintln(writer, ip)
		return err
	}
	return nil
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
