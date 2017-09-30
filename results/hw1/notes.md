## Go
Go concurrency model is based on Tony Hoare's [CSP](https://en.wikipedia.org/wiki/Communicating_sequential_processes).

In Go, Goroutines are lightweight counterparts of threads. They can be started by prepending a function call with `go` statement.

`go foo()` starts a new goroutine. If you can spare about half an hour, do have
a look at how they work and play with the code at [Tour of Go](https://tour.golang.org/concurrency/1).

For this assignment, I've used the terms "thread" and "goroutine" interchangeably to imply concurrent execution.

The goroutines are multiplexed to fewer number of OS threads. It's possible that a single thread may be handling a thousand goroutines. If a goroutine in that thread blocks, another thread is created and other goroutines are moved to new thread by Go runtime.

The difference between goroutines and threads in terms of memory consumption, cost of creation and management, context switching is well documented [here](http://blog.nindalf.com/how-goroutines-work/). The references cited by it are also worth a read and watch.

> Instead of explicitly using locks to mediate access to shared data, Go encourages the use of channels to pass references to data between goroutines. 
> 
> From [Go blog](https://blog.golang.org/share-memory-by-communicating)

### Reading Go code
Wherever possible, comments are added in client.go file to make things clear.

## Speed of Go vs Java
For 100 concurrent requests for 100 iterations, my primary results were comparable
to Java results of some friends.

Go was not inherently faster out of the box for my program. The use
of optimized http library and a powerful Desktop with 64 GB RAM have helped me
achieve 1000 times speedup over my previous Go code.

## Reported statistics and conspicuous absence of two of them:
All results are `echo`ed into files inside results directory.
Some statistics are included in the table:
| Description        | Value           |
| ------------- |:-------------:|
| Wall time      | 0.0103 s |
| Requests sent      | threads\* iterations|
| Successfull requests       | threads\* iterations|
| Mean latency | GET: 3.94111e-06 s POST: 6.32324e-06 s|
| Median latency | GET: 0.000002 s POST: 0.000003 s |
| 95 percentile latency | GET: 0.000005 s POST: 0.000007 s |
| 99 percentile latency | GET: 0.000022 s POST: 0.000024  s|

Percentiles, means, and medians have been counted via Unix coreutils like so:
```bash
❯ cat aws-fin-100-100.txt | sort -rk1 | sed '100q;d'
GET: 0.000022 POST: 0.000024

❯ awk '{ total += $2 } END { print total/NR }' aws-fin-100-100.txt
3.94111e-06

❯ awk '{ total += $4 } END { print total/NR }' aws-fin-100-100.txt
6.32324e-06
```
Awk one liners are due to [this](https://stackoverflow.com/questions/3122442/how-do-i-calculate-the-mean-of-a-column).

Total requests sent and total successful requests are not reported because _they are guaranteed if the program finishes execution normally_.

In Go, channels are a primary mechanism of communication. Sending data to and reading data from a channel blocks until the sending or receiving variable is ready. All requests times are communicated via a single channel called `result` in my program.

For this assignment, the total number of requests are no_of_threads\*no_of_iterations. In my implementation, in `main` function, I run a `for` loop for that many number of times. `main` will _not return (program will not halt) unless all the threads\*iterations values are received_.

Further, an error is reported if any of the requests are unsuccessful.

## Decision of not using Go's default http library
Initially, Go's default http library was used. For the test runs, latencies
with this library were slower or on par Java implementations of some students.
Several third party http libraries are available for Go, and I tried to test them.

For my tests, and even for production use, fasthttp library for Go is faster than Go's default library. Why this is the case and what the trade-offs are is explained [here](https://stackoverflow.com/questions/41627931/in-golang-packages-why-is-fasthttp-faster-than-net-http):

## Stress tests and breaking things
Since creating goroutines is cheaper than a thread, I can easily create millions of goroutines in hundreds of microseconds without risking out of memory exception.
See it in action [here](https://play.golang.org/p/sLRTI6p3ie)

However, if we use Go's default http library to make GET and POST requests, we receive a response (a Unix [file](https://en.wikipedia.org/wiki/Everything_is_a_file)) which must be closed after we deal with errors.

Every OS has some default value of number of files that can be open at any time. For Linux and OS X machines I used, this value seems to be 1024. If I create 1000 threads, I receive an errors that too many files are open.

Fortunately, this default value can be changed. Even better, Go's fastHttp library has a different implementation which avoids this issue altogether (I still need to understand how).

If that library is used, creating 10 million goroutines (10 million concurrent GET and POST requests for 100 iterations) is easy on a 64 GB machine, like so :

As you see in the prompt in the screenshot, the previous command took 3 minutes and
16 seconds to execute. It utilized 1 million concurrent requests for 100 iterations.

If 10 million requests are completed within half an hour, they will be reported before
deadline. (Edit They couldn't get completed. They made the 64 GB RAM, 512 GB SSD desktop run out of memory. Screenshots of such breakdowns are in hw1 folder of this repo.)

The fact that server stayed alive after so many requests seems a bit unreasonable to me even when I am using Go.
