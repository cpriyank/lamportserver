package lamportserver

import (
	"fmt"
	"github.com/buaazp/fasthttprouter"
	"github.com/valyala/fasthttp"
	"log"
	"os"
	"time"
)

var getTrigger = make(chan bool, numStats)
var getResponseLogChan = make(chan *LatencyStat, numStats)
var postResponseLogChan = make(chan *LatencyStat, numStats)
var dbGETLatencyLogChan = make(chan *LatencyStat, numStats)
var dbPOSTLatencyLogChan = make(chan *LatencyStat, numStats)
var receiveTrigger = make(chan bool, numStats)

type LatencyStat struct {
	Latency   float64
	TimeStamp int64
}

// A middleware which logs response time of a request handler given
func logHandlers(h fasthttp.RequestHandler, endpoint string) fasthttp.RequestHandler {
	return fasthttp.RequestHandler(func(ctx *fasthttp.RequestCtx) {
		start := time.Now()
		h(ctx)
		latency := time.Since(start).Seconds()
		// TODO: poor if block
		if endpoint == "POST" {
			postResponseLogChan <- &LatencyStat{latency, time.Now().UnixNano()}
		} else {
			getResponseLogChan <- &LatencyStat{latency, time.Now().UnixNano()}
		}
	})
}

// MultiParams is the multi params handler
func vertStats(ctx *fasthttp.RequestCtx) {
	getTrigger <- true
	skierID, dayNum := parseQuery(ctx)
	start := time.Now()
	verticals, lifts := queryDB(skierID, dayNum)
	dbGetLatency := time.Since(start).Seconds()
	dbGETLatencyLogChan <- &LatencyStat{dbGetLatency, time.Now().UnixNano()}
	fmt.Fprintf(ctx, "%s%s", verticals, lifts)

}

// MultiParams is the multi params handler
func loadStats(ctx *fasthttp.RequestCtx) {
	// fmt.Fprintf(ctx, "hi, %s, %s %s %s %s!\n", ctx.UserValue("resortID"), ctx.UserValue("dayNum"), ctx.UserValue("skierID"), ctx.UserValue("liftID"), ctx.UserValue("timeStamp"))
	receiveTrigger <- true

	fmt.Fprintf(ctx, "hi")
	stat, err := parse(ctx)
	go func() {
		// fmt.Println(stat)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%v", err)
		}
		stat.load()
	}()
}

func Serve() {
	router := fasthttprouter.New()
	router.GET("/myvert/:skierID/:dayNum", logHandlers(vertStats, "GET"))
	router.POST("/load/:resortID/:dayNum/:skierID/:liftID/:timeStamp", logHandlers(loadStats, "POST"))

	go writeToDB()
	go fanInLatencies()

	log.Fatal(fasthttp.ListenAndServe(":8000", router.Handler))
}
