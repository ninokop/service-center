package metrics

import (
	"bytes"
	"fmt"
	"io/ioutil"
	text_template "text/template"
	"time"

	// "github.com/ServiceComb/service-center/pkg/metrics/template"
	"github.com/ServiceComb/service-center/pkg/util"
	pm "github.com/prometheus/client_golang/prometheus"
	dto "github.com/prometheus/client_model/go"
)

func init() {
	go startReportMetrics()
}

func startReportMetrics() {
	interval := time.Duration(10) * time.Second
	tick := time.NewTicker(interval)
	defer tick.Stop()

	for {
		select {
		case <-tick.C:
			m, err := metricsCollect()
			if err != nil {
				fmt.Printf("########## Error occurred during metrics collection: %v\n", err)
				continue
			}
			m.WriteToFiles()
			// case <-done:
			// 	return
		}
	}
}

var targetMetrics = util.NewStringSets("go_threads", "http_request_total", "process_cpu_seconds_total", "process_resident_memory_bytes")
var baseMetricMap = map[string]float64{}

func metricsCollect() (*Metrics, error) {
	mfs, err := pm.DefaultGatherer.Gather()
	if err != nil {
		return nil, err
	}
	for _, mf := range mfs {
		if targetMetrics.Has(mf.GetName()) {
			baseMetricMap[mf.GetName()] = valueOfMetric(mf)
		}
	}
	return &Metrics{BaseMetricMap: baseMetricMap}, nil
}

func valueOfMetric(mf *dto.MetricFamily) float64 {
	if len(mf.GetMetric()) == 0 {
		return 0
	}
	switch mf.GetType() {
	case dto.MetricType_COUNTER:
		return *mf.GetMetric()[0].GetCounter().Value
	case dto.MetricType_GAUGE:
		return *mf.GetMetric()[0].GetGauge().Value
	default:
		return 0
	}
}

type Metrics struct {
	BaseMetricMap map[string]float64
	// CPU         string
	// Memery      string
	// Threads     string

	// HttpRequests string
	// // SuccessfulRate
	// // Lantency
	// // Qps
}

type CseMetrics struct {
	PluginID      string
	MetricsSource string
	ScopeName     string
	TimeStamp     string
	InfaceName    string
	MetricsKey    string
	MetricsValue  float64
}

func (ms *Metrics) WriteToFiles() {
	fmt.Println("########write to files######", ms)

	tmpl := text_template.New("metrics.dat")
	t, err := tmpl.ParseFiles("./pkg/metrics/template/metrics.dat")
	if err != nil {
		fmt.Printf("Parse template file error: %v\n", err)
		return
	}

	for k, v := range ms.BaseMetricMap {
		m := CseMetrics{
			PluginID:      "nino",
			MetricsSource: "10.120.195.2",
			ScopeName:     "sc",
			InfaceName:    "inface_name",
			MetricsKey:    k,
			MetricsValue:  v,
		}

		buffer := new(bytes.Buffer)
		if err := t.Execute(buffer, m); err != nil {
			fmt.Printf("Execute get template error: %v\n", err)
			return
		}

		outPath := fmt.Sprintf("./pkg/metrics/output/metrics_%s.dat", k)
		err = ioutil.WriteFile(outPath, buffer.Bytes(), 0644)
		if err != nil {
			fmt.Printf("Content Files is: %s\n", string(buffer.Bytes()))
			return
		}
	}
}
