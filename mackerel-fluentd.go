package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	mp "github.com/mackerelio/go-mackerel-plugin"
)

type FluentdMetrics struct {
	Target   string
	Tempfile string

	plugins []FluentdPluginMetrics
	err     error
}

type FluentdPluginMetrics struct {
	RetryCount     uint64            `json:"retry_count"`
	OutputPlugin   bool              `json:"output_plugin"`
	Config         map[string]string `json:"config"`
	Type           string            `json:"type"`
	PluginCategory string            `json:"plugin_category"`
	PluginID       string            `json:"plugin_id"`
}

type FluentMonitorJSON struct {
	Plugins []FluentdPluginMetrics `json:"plugins"`
}

func (f FluentdMetrics) fetch() ([]FluentdPluginMetrics, error) {
	return []FluentdPluginMetrics{}, nil
}

func (f *FluentdMetrics) Init() {
	resp, err := http.Get(f.Target)
	if err != nil {
		f.err = err
		return
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		f.err = err
		return
	}
	var j FluentMonitorJSON
	err = json.Unmarshal(body, &j)
	f.plugins = j.Plugins
	f.err = err
}

func (f FluentdMetrics) FetchMetrics() (map[string]float64, error) {
	metrics := map[string]float64{}
	for _, plugin := range f.plugins {
		metrics[plugin.PluginID] = float64(plugin.RetryCount)
	}
	return metrics, f.err
}

func (f FluentdMetrics) GraphDefinition() map[string](mp.Graphs) {
	metrics := [](mp.Metrics){}
	for _, p := range f.plugins {
		metrics = append(metrics, mp.Metrics{
			Name:  p.PluginID,
			Label: p.PluginID,
			Diff:  false})
	}

	return map[string](mp.Graphs){
		"fluentd.retry": mp.Graphs{
			Label:   "Fluentd retry count",
			Unit:    "integer",
			Metrics: metrics,
		},
	}
}

func main() {
	host := flag.String("host", "localhost", "fluentd monitor_agent port")
	port := flag.String("port", "24220", "fluentd monitor_agent port")
	tempFile := flag.String("tempfile", "", "Temp file name")
	flag.Parse()

	f := FluentdMetrics{
		Target:   fmt.Sprintf("http://%s:%s/api/plugins.json", *host, *port),
		Tempfile: *tempFile,
	}
	f.Init()
	helper := mp.NewMackerelPlugin(f)

	if *tempFile != "" {
		helper.Tempfile = *tempFile
	} else {
		helper.Tempfile = fmt.Sprintf("/tmp/mackerel-plugin-fluentd-%s-%s", *host, *port)
	}

	if os.Getenv("MACKEREL_AGENT_PLUGIN_META") != "" {
		helper.OutputDefinitions()
	} else {
		helper.OutputValues()
	}
}
