package lamportserver

import (
	"fmt"
	"time"
)

func tickClock(delay time.Duration, iterations int) {
	for {
		for i := 0; i < iterations; i++ {
			fmt.Printf("\r%d", i)
			time.Sleep(delay)
		}
	}
}
