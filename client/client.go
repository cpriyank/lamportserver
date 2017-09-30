package main

import (
	"fmt"
	"github.com/valyala/fasthttp"
	"os"
	"strconv"
	"time"
)

// Function makeRequest makes an http GET and POST requests to a URI, and
// reports errors if any. Then, it saves the response time of each of the calls
// in a channel (call it "a synchronized variable called result")
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

	// Assign the time taken by requests in a channel passed to this goroutine.
	result <- fmt.Sprintf("GET: %f POST: %f", getSecs, postSecs)
}

// MakeRequestHelper iteratively makes requests to the URI for specified number
// of iterations.
func MakeRequestHelper(url string, result chan<- string, iterations int) {
	for i := 0; i < iterations; i++ {
		makeRequest(url, result)
	}
}

func main() {
	args := os.Args[1:]
	// default args
	threads := 100
	iterations := 100
	portString := "8000"
	url := "http://localhost" // Not giving AWS url to prevent DDOS.
	buffers := 1

	// parse the arguments. Assuming user will pass all 5 args
	// or none at all
	if len(args) != 0 {
		threadString := args[0]
		iterationString := args[1]
		url = args[2]
		bufferString := args[4]
		portString = args[3]

		threadsInt, err := strconv.Atoi(threadString)
		threads = threadsInt
		if err != nil {
			fmt.Printf("%v", err)
		}
		iterations, err = strconv.Atoi(iterationString)
		if err != nil {
			fmt.Printf("%v", err)
		}
		// I have additional argument which is the length of buffer
		// in go channel (Go concept. Read more here on it on Tour of Go), however,
		// after running a few tests, I didn't notice any difference in latencies
		// whether I use buffered channel or not.
		buffers, err = strconv.Atoi(bufferString)
		if err != nil {
			fmt.Printf("%v", err)
		}

	}

	url = fmt.Sprintf("%s:%s", url, portString)

	// result channel is the safe and only way for threads to communicate
	// with each other in this program.
	result := make(chan string, buffers)
	resultSlice := make([]string, threads*iterations)
	start := time.Now()

	for i := 0; i < threads; i++ {
		// go statement calls the function that succeeds in a background.
		// It's like running a process in background in a Unix shell.
		// An example: https://tour.golang.org/concurrency/1
		go MakeRequestHelper(url, result, iterations)
	}

	for i := 0; i < threads*iterations; i++ {
		// as the results arrive, save them in an array. Note that
		// reading from result channel blocks until it has a non empty
		// value. By running this for loop threads * iterations times,
		// it is guaranteed that execution will only proceed if/when
		// all the call results complete
		resultSlice[i] = <-result
	}
	fmt.Println("All threads took in total:", time.Since(start).Seconds())
	fmt.Println("Time taken by each of the threads for GET and POST:")
	for _, singleResult := range resultSlice {
		fmt.Println(singleResult)
	}
}
