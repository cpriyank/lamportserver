package lamportserver

import (
	"bytes"
	"github.com/wcharczuk/go-chart"
)

func separateFields(stats []*LatencyStat) ([]float64, []float64) {
	latencies := make([]float64, len(stats))
	timestamps := make([]float64, len(stats))
	for i := range stats {
		latencies[i] = stats[i].Latency
		timestamps[i] = float64(stats[i].TimeStamp)
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
