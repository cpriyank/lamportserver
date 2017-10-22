package lamportserver

import (
	"github.com/valyala/fasthttp"
	"strconv"
)

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
	stat := &skierStat{resortID, dayNum, skierID, liftID, timeStamp}
	return stat, nil
}
