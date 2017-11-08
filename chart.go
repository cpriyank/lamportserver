package lamportserver

import (
	"bytes"
	"fmt"
	"github.com/wcharczuk/go-chart"
	"os"
)

func separateFields(stats []*LatencyStat) ([]float64, []float64) {
	latencies := make([]float64, len(stats))
	timestamps := make([]float64, len(stats))
	for i := range stats {
		latencies[i] = stats[i].Latency
		timestamps[i] = float64(stats[i].TimeStamp)
		fmt.Println(latencies[i], timestamps[i])
	}
	return timestamps, latencies
}

// TODO: currently nameToSave must have PNG extension
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
	file, err := os.Create(nameToSave)
	defer file.Close()
	_, err = file.Write(buffer.Bytes())
	return err
}
