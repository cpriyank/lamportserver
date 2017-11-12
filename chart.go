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
	}
	return timestamps, latencies
}

// TODO: currently nameToSave must have PNG extension
func chartStat(stats []*LatencyStat, nameToSave string) error {

	timestamps, latencies := separateFields(stats)
	graph := chart.Chart{
		XAxis: chart.XAxis{
			Name:      "Time stamps",
			NameStyle: chart.StyleShow(),
			Style:     chart.StyleShow(),
		},
		YAxis: chart.YAxis{
			Name:      "latencies",
			NameStyle: chart.StyleShow(),
			Style:     chart.StyleShow(),
			ValueFormatter: func(v interface{}) string {
				if vf, isFloat := v.(float64); isFloat {
					return fmt.Sprintf("%0.8f", vf)
				}
				return ""
			},
		},
		Series: []chart.Series{
			chart.ContinuousSeries{
				Style: chart.Style{
					Show:        true,
					StrokeColor: chart.GetDefaultColor(0).WithAlpha(64),
					FillColor:   chart.GetDefaultColor(0).WithAlpha(64),
				},
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
