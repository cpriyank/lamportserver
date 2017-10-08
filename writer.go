package lamportserver

import (
	"github.com/gocql/gocql"
)

type skierStat struct {
	resortID  int
	dayNum    int
	skierID   int
	liftID    int
	timeStamp int
}

var statCache []*skierStat = make([]*skierStat, numStats)
var statChan = make(chan *skierStat, concurrency)

func writeUsingStatChan() {
	for i := 0; i < numStats; i++ {
		statCache[i] = <-statChan
	}

}

func (stat *skierStat) load() {
	statChan <- stat
}
