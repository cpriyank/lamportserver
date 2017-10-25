package lamportserver

import (
	"fmt"
	"github.com/buaazp/fasthttprouter"
	"github.com/valyala/fasthttp"
	"log"
	"os"
)

var getTrigger = make(chan bool, numStats)
var resultsOfGet = make(chan string, numSkiers)

// MultiParams is the multi params handler
func vertStats(ctx *fasthttp.RequestCtx) {
	getTrigger <- true
	skierID, dayNum := parseQuery(ctx)
	verticals, lifts := queryDB(skierID, dayNum)
	fmt.Fprintf(ctx, "%s%s", verticals, lifts)

}

// QueryArgs is used for uri query args test #11:
// if the req uri is /ping?name=foo, output: Pong! foo
// if the req uri is /piNg?name=foo, redirect to /ping, output: Pong!
// func QueryArgs(ctx *fasthttp.RequestCtx) {
// 	name := ctx.QueryArgs().Peek("stats")
// 	fmt.Fprintf(ctx, "Pong! %s\n", string(stats))
// }

var receiveTrigger = make(chan bool, numStats)

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
	router.GET("/myvert/:skierID/:dayNum", vertStats)
	router.POST("/load/:resortID/:dayNum/:skierID/:liftID/:timeStamp", loadStats)
	// router.POST("/load", QueryArgs)
	// go writeUsingStatChan()
	// for i := 0; i < dbConnPoolSize; i++ {
	go writeToDB()

	log.Fatal(fasthttp.ListenAndServe(":8000", router.Handler))
}
