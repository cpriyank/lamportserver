package lamportserver

import (
	"fmt"
	"time"
)

const (
	numStats       = 800000
	numSkiers      = 40000
	concurrency    = 100
	dbConnPoolSize = 50
)

// var statCache = make([]*skierStat, numStats)
var statChan = make(chan *skierStat, numStats)
var getResponseLatencies = make([]*LatencyStat, numSkiers)
var postResponseLatencies = make([]*LatencyStat, numStats)
var dbWriteLatencies = make([]*LatencyStat, numStats)
var dbReadLatencies = make([]*LatencyStat, numSkiers)

func fanInLatencies() {
	for {
		select {
		case responseStat := <-getResponseLogChan:
			go passLatencyToMQ(responseStat, "getCTX")
		case responseStat := <-postResponseLogChan:
			go passLatencyToMQ(responseStat, "postCTX")
		case responseStat := <-dbGETLatencyLogChan:
			go passLatencyToMQ(responseStat, "getDB")
		case responseStat := <-dbPOSTLatencyLogChan:
			go passLatencyToMQ(responseStat, "postDB")
		}
	}
}

func closeChans() {
	if getCounter == numSkiers {
		close(getTrigger)
		close(dbGETLatencyLogChan)
		close(getResponseLogChan)
	}
	if postCounter == numSkiers {
		close(receiveTrigger)
		close(dbPOSTLatencyLogChan)
		close(postResponseLogChan)
	}
}

func passLatencyToMQ(stat *LatencyStat, classification string) {
	switch classification {
	case "getCTX":
		getResponseLatencies = append(getResponseLatencies, stat)
		if len(getResponseLatencies) == numSkiers {
			fileName := fmt.Sprintf("get-%d.PNG", time.Now().UnixNano())
			chartStat(getResponseLatencies, fileName)
		}
	case "postCTX":
		postResponseLatencies = append(postResponseLatencies, stat)
		if len(postResponseLatencies) == numStats {
			fileName := fmt.Sprintf("post-%d.PNG", time.Now().UnixNano())
			chartStat(postResponseLatencies, fileName)
		}
	case "getDB":
		dbReadLatencies = append(dbReadLatencies, stat)
		if len(dbReadLatencies) == numSkiers {
			fileName := fmt.Sprintf("dbread-%d.PNG", time.Now().UnixNano())
			chartStat(dbReadLatencies, fileName)
		}
	case "postDB":
		dbWriteLatencies = append(dbWriteLatencies, stat)
		if len(dbWriteLatencies) == numStats {
			fileName := fmt.Sprintf("dbwrite-%d.PNG", time.Now().UnixNano())
			chartStat(dbWriteLatencies, fileName)
		}
	}
}

func (stat *skierStat) load() {
	statChan <- stat
}
