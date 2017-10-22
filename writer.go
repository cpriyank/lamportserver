package lamportserver

// import (
// 	"github.com/gocql/gocql"
// )

type skierStat struct {
	resortID  interface{}
	dayNum    interface{}
	skierID   interface{}
	liftID    interface{}
	timeStamp interface{}
}

var statCache = make([]*skierStat, numStats)
var statChan = make(chan *skierStat, concurrency)

func writeUsingStatChan() {
	for i := 0; i < numStats; i++ {
		statCache[i] = <-statChan
	}

}

func (stat *skierStat) load() {
	statChan <- stat
}
