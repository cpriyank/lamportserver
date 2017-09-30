package main

import (
	"fmt"
	"github.com/valyala/fasthttp"
	"os"
	"strconv"
	"time"
)

func makeRequest(url string, result chan<- string) {
	start := time.Now()
	statusCode, body, err := fasthttp.Get(nil, url)
	// fasthttp.Get(nil, url)
	if err != nil {
		fmt.Printf("%v", err)
	}
	if statusCode != fasthttp.StatusOK {
		fmt.Fprintf(os.Stderr, "Unexpected status code: %d. Expecting %d", statusCode, fasthttp.StatusOK)
	}
	// this is hard-coded
	if len(body) != 27 {
		fmt.Printf("%v", err)
	}
	getSecs := time.Since(start).Seconds()

	statusCode, body, err = fasthttp.Post(nil, url, nil)
	// fasthttp.Post(nil, url, nil)
	if err != nil {
		fmt.Printf("%v", err)
	}
	if statusCode != fasthttp.StatusOK {
		fmt.Fprintf(os.Stderr, "Unexpected status code: %d. Expecting %d", statusCode, fasthttp.StatusOK)
	}
	// this is hard-coded
	if len(body) != 27 {
		fmt.Printf("%v", err)
	}
	postSecs := time.Since(start).Seconds()

	result <- fmt.Sprintf("GET: %f POST: %f, at %v", getSecs, postSecs, time.Now().String())
}

func MakeRequestHelper(url string, result chan<- string, iterations int) {
	for i := 0; i < iterations; i++ {
		makeRequest(url, result)
	}
}

func main() {
	args := os.Args[1:]
	threadString := args[0]
	iterationString := args[1]
	url := args[2]
	bufferString := args[4]
	portString := args[3]
	// port
	threads, err := strconv.Atoi(threadString)
	if err != nil {
		fmt.Printf("%v", err)
	}
	iterations, err := strconv.Atoi(iterationString)
	if err != nil {
		fmt.Printf("%v", err)
	}
	buffers, err := strconv.Atoi(bufferString)
	if err != nil {
		fmt.Printf("%v", err)
	}
	// port, err := strconv.Atoi(portString)
	// if err != nil {
	// fmt.Printf("%v", err)
	// }
	url = fmt.Sprintf("%s:%s", url, portString)
	buffers = 1
	result := make(chan string, buffers)
	resultSlice := make([]string, threads*iterations)
	start := time.Now()
	for i := 0; i < threads; i++ {
		go MakeRequestHelper(url, result, iterations)
	}
	for i := 0; i < threads*iterations; i++ {
		resultSlice[i] = <-result
	}
	fmt.Println("All threads took in total:", time.Since(start).Seconds())
	fmt.Println("Time taken by each of the threads for GET and POST:")
	for _, singleResult := range resultSlice {
		fmt.Println(singleResult)
	}
}
