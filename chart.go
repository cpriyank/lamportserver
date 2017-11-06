package lamportserver

import (
	"bytes"
	"github.com/wcharczuk/go-chart"
)

type LatencyStat struct {
	Latency   float64
	TimeStamp float64
}

func separateFields(stats []*LatencyStat) ([]float64, []float64) {
	latencies := make([]float64, len(stats))
	timestamps := make([]float64, len(stats))
	for i := range stats {
		latencies[i] = stats[i].Latency
		timestamps[i] = stats[i].TimeStamp
	}
	return timestamps, latencies
}

func chartStat(stats []*LatencyStat, nameToSave string) error {

	timestamps, latencies := separateFields(stats)
	graph := chart.Chart{
		Series: []chart.Series{
			chart.ContinuousSeries{
				XValues: timestamps,
				YValues: latencies,
			},
		},
	}
	buffer := bytes.NewBuffer([]byte{})
	err := graph.Render(chart.PNG, buffer)
	return err
}
