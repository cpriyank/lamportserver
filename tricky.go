package main

import (
	"fmt"
	"os"
	"strconv"
	"time"
)

func main() {
	args := os.Args[1:]
	threads := 10
	iterations := 100
	ipAddr := "aws default ip"
	port := 8000

	if len(args) > 0 {
		threads, iterations, ipAddr, port, err := parseArgs(args)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error parsing command line args: %v\n", err)
			os.Exit(1)
		}
	}

	logs := make([][]Time, iterations)
	for i := range iterations {
		logs[i] = make([]Time, threads)
	}

}

func timeTrack(start time.Time, name string) {
	elapsed := time.Since(start)
	fmt.Printf("%s took %s", name, elapsed)
}

func parseArgs(args) (int, int, string, int, error) {
	threads, err := strconv.Atoi(args[0])
	if err != nil {
		return -1, -1, "", -1, fmt.Printf("Error parsing %d, %v", args, err)
	}

	iterations, err := strconv.Atoi(args[1])
	if err != nil {
		return -1, -1, "", -1, fmt.Printf("Error parsing %d, %v", args, err)
	}

	ip = args[2]
	port, err := strconv.Atoi(args[3])
	if err != nil {
		return -1, -1, "", -1, fmt.Printf("Error parsing %d, %v", args, err)
	}
	return threads, iterations, ip, port, err

}
