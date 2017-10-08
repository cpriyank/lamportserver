package lamportserver

import (
	"fmt"
	"github.com/buaazp/fasthttprouter"
	"github.com/valyala/fasthttp"
	"log"
)

const (
	concurrency = 100
	numStats    = 999
)

// MultiParams is the multi params handler
func vertStats(ctx *fasthttp.RequestCtx) {
	fmt.Fprintf(ctx, "hi, %s, %s!\n", ctx.UserValue("skierID"), ctx.UserValue("dayNum"))
}

// QueryArgs is used for uri query args test #11:
// if the req uri is /ping?name=foo, output: Pong! foo
// if the req uri is /piNg?name=foo, redirect to /ping, output: Pong!
// func QueryArgs(ctx *fasthttp.RequestCtx) {
// 	name := ctx.QueryArgs().Peek("stats")
// 	fmt.Fprintf(ctx, "Pong! %s\n", string(stats))
// }

// MultiParams is the multi params handler
func loadStats(ctx *fasthttp.RequestCtx, incoming chan<- *skierStat) {
	// fmt.Fprintf(ctx, "hi, %s, %s %s %s %s!\n", ctx.UserValue("resortID"), ctx.UserValue("dayNum"), ctx.UserValue("skierID"), ctx.UserValue("liftID"), ctx.UserValue("timeStamp"))
	fmt.Fprintf(ctx, "hi")
	incoming := &skierStat{resortID: ctx.UserValue("resortID"), dayNum: ctx.UserValue("dayNum"), skierID: ctx.UserValue("skierID"), liftID: ctx.UserValue("liftID"), timeStamp: ctx.UserValue("timeStamp")}
	incoming.load()
}

func Serve() {
	router := fasthttprouter.New()
	router.GET("/myvert/:skierID/:dayNum", vertStats)
	router.POST("/load/:resortID/:dayNum/:skierID/:liftID/:timeStamp", loadStats)
	// router.POST("/load", QueryArgs)
	go func() {
		writeUsingStatChan()
	}()
	log.Fatal(fasthttp.ListenAndServe(":8000", router.Handler))
}
