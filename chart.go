package lamportserver

import (
	"bytes"
	"fmt"
	"github.com/wcharczuk/go-chart"
	"os"
	"sort"
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

func (s LatencyStat) String() string {
	return fmt.Sprintf("Latency:%f, TimeStamp: %d", s.Latency, s.TimeStamp)
}


func averageLatency(stats []*LatencyStat) (float64, error) {
	if len(stats) == 0 {
		return 0, fmt.Errorf("input sequence is empty")
	}
	var sum float64
	for i := range stats {
		sum += stats[i].Latency
	}
	return sum / float64(len(stats)), nil

}

func medianAndP9s(stats []*LatencyStat) (float64, float64, float64) {
	sort.Sort(ByLatency(stats))
	medianIndex := len(stats) / 2
	p95Index := int(float64(len(stats)) * 0.95)
	p99Index := int(float64(len(stats)) * 0.99)
	return stats[medianIndex].Latency, stats[p95Index].Latency, stats[p99Index].Latency
}


type ByLatency []*LatencyStat

func (s ByLatency) Len() int {
	return len(s)
}

func (s ByLatency) Less(i, j int) bool {
	return s[i].Latency < s[j].Latency
}

func (s ByLatency) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}


// TODO: currently nameToSave must have PNG extension
func chartStat(stats []*LatencyStat, nameToSave string) error {

	avgGET, err := averageLatency(stats)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v", err)
	}
	medianGET, p95GET, p99GET := medianAndP9s(stats)
	fmt.Printf("%s avg = %f, median = %f, p95 = %f, p99 = %f\n",nameToSave, avgGET, medianGET, p95GET, p99GET)

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
	err = graph.Render(chart.PNG, buffer)
	file, err := os.Create(nameToSave)
	defer file.Close()
	_, err = file.Write(buffer.Bytes())
	return err
}
