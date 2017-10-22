package lamportserver

import (
	"fmt"
)

type skierStat struct {
	resortID  int
	dayNum    int
	skierID   int
	liftID    int
	timeStamp int
}

var statCache = make([]*skierStat, numStats)
var statChan = make(chan *skierStat, concurrency)

func writeUsingStatChan() {
	for i := 0; i < numStats; i++ {
		statCache[i] = <-statChan
	}
	for _, stat := range statCache {
		fmt.Println(stat)
	}

}

func (stat *skierStat) load() {
	statChan <- stat
}
