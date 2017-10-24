package lamportserver

import (
	"fmt"
	"time"
)

const (
	numStats       = 800000
	concurrency    = 100
	dbConnPoolSize = 100
)

var statCache = make([]*skierStat, numStats)
var statChan = make(chan *skierStat, concurrency)

func writeUsingStatChan() {
	start := time.Now()
	for i := 0; i < numStats; i++ {
		statCache[i] = <-statChan
	}
	// for i, stat := range statCache {
	fmt.Println("took", time.Since(start))
	// }

	// mapSkierToDaysToLiftID(statCache)
	writeToDB()

}

func (stat *skierStat) load() {
	statChan <- stat
}
