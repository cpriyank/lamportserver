package lamportserver

const (
	numStats       = 800000
	concurrency    = 100
	dbConnPoolSize = 50
)

// var statCache = make([]*skierStat, numStats)
var statChan = make(chan *skierStat, numStats)

func fanInLatencies() {
	for {
		select {
		case responseStat := <-getResponseLogChan:
			passLatencyToMQ(responseStat, "getCTX")
		case responseStat := <-postResponseLogChan:
			passLatencyToMQ(responseStat, "postCTX")
		case responseStat := <-dbGETLatencyLogChan:
			passLatencyToMQ(responseStat, "getDB")
		case responseStat := <-dbPOSTLatencyLogChan:
			passLatencyToMQ(responseStat, "postDB")
		}
	}
}

func passLatencyToMQ(stat *LatencyStat, classification string) {
}

func (stat *skierStat) load() {
	statChan <- stat
}
