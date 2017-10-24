package lamportserver

import (
	"github.com/valyala/fasthttp"
	"strconv"
)

const (
	maxInt = 1<<31 - 1
)

func liftIDToVertical(liftID int) int {
	switch {
	case liftID <= 10:
		return 200
	case liftID <= 20:
		return 300
	case liftID <= 30:
		return 400
	case liftID <= 40:
		return 500
	}
	return maxInt
}

func parse(ctx *fasthttp.RequestCtx) (*skierStat, error) {
	resortIDString := ctx.UserValue("resortID").(string)
	resortID, err := strconv.Atoi(resortIDString)
	if err != nil {
		return nil, err
	}
	dayNumString := ctx.UserValue("dayNum").(string)
	dayNum, err := strconv.Atoi(dayNumString)
	if err != nil {
		return nil, err
	}
	skierIDString := ctx.UserValue("skierID").(string)
	skierID, err := strconv.Atoi(skierIDString)
	if err != nil {
		return nil, err
	}
	liftIDString := ctx.UserValue("liftID").(string)
	liftID, err := strconv.Atoi(liftIDString)
	if err != nil {
		return nil, err
	}
	timeStampString := ctx.UserValue("timeStamp").(string)
	timeStamp, err := strconv.Atoi(timeStampString)
	if err != nil {
		return nil, err
	}
	verticals := liftIDToVertical(liftID)
	stat := &skierStat{resortID, dayNum, skierID, liftID, timeStamp, verticals}
	return stat, nil
}
