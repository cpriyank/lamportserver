package lamportserver

const (
	numStats    = 800000
	concurrency = 100
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
	// for i, stat := range statCache {
	// 	fmt.Println(i, stat)
	// }

	// mapSkierToDaysToLiftID(statCache)

}

func (stat *skierStat) load() {
	statChan <- stat
}
