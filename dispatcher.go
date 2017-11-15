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
var postResponseLatencies []*LatencyStat
var dbWriteLatencies []*LatencyStat
var dbReadLatencies []*LatencyStat

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
		close(dbGETLatencyLogChan)
		close(getResponseLogChan)
	}
	if dbPOSTCounter == numStats {
		close(receiveTrigger)
		close(dbPOSTLatencyLogChan)
		close(postResponseLogChan)
	}
}

var getthres = []int{1000, 2000, 3000, 4000, 5000, 6000, 7000, 8000, 9000, 1000}

var postres = []int{10000, 20000, 30000, 40000, 50000, 60000, 70000, 80000, 90000, 100000, 110000, 120000, 130000, 140000, 150000, 160000, 170000, 180000, 190000, 200000}

func findinIntArr(arr []int, val int) bool {
	for _, v := range arr {
		if v == val {
		return true
		}
	}
	return false
}

func passLatencyToMQ(stat *LatencyStat, classification string) {
	switch classification {
	case "getCTX":
		getResponseLatencies = append(getResponseLatencies, stat)
		if findinIntArr(getthres, len(getResponseLatencies)) {
			fmt.Println("get response")
			fileName := fmt.Sprintf("get-%d.PNG", time.Now().UnixNano())
			chartStat(getResponseLatencies, fileName)
		}
	case "postCTX":
		postResponseLatencies = append(postResponseLatencies, stat)
		if findinIntArr(postres, len(postResponseLatencies)) {
			fmt.Println("post response")
			fileName := fmt.Sprintf("post-%d.PNG", time.Now().UnixNano())
			chartStat(postResponseLatencies, fileName)
		}
	case "getDB":
		dbReadLatencies = append(dbReadLatencies, stat)
		if findinIntArr(getthres, len(dbReadLatencies)) {
			fmt.Println("dbread response")
			fileName := fmt.Sprintf("dbread-%d.PNG", time.Now().UnixNano())
			chartStat(dbReadLatencies, fileName)
		}
	case "postDB":
		dbWriteLatencies = append(dbWriteLatencies, stat)
		if findinIntArr(postres, len(dbWriteLatencies)){
			fmt.Println("dbwrite response")
			fileName := fmt.Sprintf("dbwrite-%d.PNG", time.Now().UnixNano())
			chartStat(dbWriteLatencies, fileName)
		}
	}
}

func (stat *skierStat) load() {
	statChan <- stat
}
