# cidrex

`cidrex` is a command-line utility for expanding IP addresses and CIDR ranges. It reads IP addresses and CIDR ranges from a file or stdin, and outputs all individual IP addresses contained within them. The utility supports both IPv4 and IPv6 addresses.

![Example output](https://i.imgur.com/YkmgTc4.png)

## Features

- Supports reading from a specified file or standard input (stdin).
- Expands CIDR ranges into individual IP addresses.
- Supports filtering only IPv4, only IPv6, or both types of addresses.

## Installation

To install `cidrex`, you need to have [Go](https://golang.org) installed on your system. You can then install the binary using the following commands:

```bash
go install github.com/d3mondev/cidrex@latest
```

## Usage

```bash
cidrex [OPTIONS] [filename]
```

If no filename is provided, cidrex reads from stdin.

### Options

* `-4, --ipv4`: Print only IPv4 addresses
* `-6, --ipv6`: Print only IPv6 addresses
* `-h, --help`: Display the help message

### Examples


1. Expand IP addresses and CIDR ranges from a file:

```bash
cidrex input.txt
```
2. Expand only IPv4 addresses from a file:

```bash
cidrex -4 input.txt
```

3. Expand only IPv6 addresses from stdin:

```
cat input.txt | cidrex -6
```

### Input Format

The input should contain one IP address or CIDR range per line. For example:

```
192.168.1.1
10.0.0.0/24
2001:db8::1
2001:db8::/120
```

### Output

The program outputs one IP address per line. For individual IP addresses, it simply outputs the address as-is. For CIDR ranges, it expands the range and outputs each individual IP address within that range.

If the program encounters any errors (e.g., invalid IP addresses or CIDR ranges), it will print error messages to stderr and continue processing the remaining input.
