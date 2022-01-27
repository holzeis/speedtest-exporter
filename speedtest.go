package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"math"
	"net/http"
	"os/exec"
	"strings"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

const (
	namespace = "speedtest"
)

var (
	ping = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "", "ping"),
		"Latency (ms)",
		nil, nil,
	)
	download = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "", "download"),
		"Download bandwidth (Mbps).",
		nil, nil,
	)
	upload = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "", "upload"),
		"Upload bandwidth (Mbps).",
		nil, nil,
	)
)

type Metrics struct {
	Download json.Number `json:"download"`
	Upload   json.Number `json:"upload"`
	Ping     json.Number `json:"ping"`
}

// Exporter collects Speedtest stats from the given server and exports them using
// the prometheus metrics package.
type Exporter struct {
}

// Describe describes all the metrics ever exported by the Speedtest exporter.
// It implements prometheus.Collector.
func (e *Exporter) Describe(ch chan<- *prometheus.Desc) {
	ch <- ping
	ch <- download
	ch <- upload
}

// Collect fetches the stats from configured Speedtest location and delivers them as Prometheus metrics.
// It implements prometheus.Collector.
func (e *Exporter) Collect(ch chan<- prometheus.Metric) {
	log.Printf("Speedtest exporter starting")

	metrics, err := e.NetworkMetrics()
	if err != nil {
		log.Printf("Error: %v", err.Error())
		return
	}

	m_download, _ := metrics.Download.Float64()
	m_upload, _ := metrics.Upload.Float64()
	m_ping, _ := metrics.Ping.Float64()

	ch <- prometheus.MustNewConstMetric(ping, prometheus.GaugeValue, math.Round(m_ping))
	ch <- prometheus.MustNewConstMetric(download, prometheus.GaugeValue, math.Round(m_download/1024/1024))
	ch <- prometheus.MustNewConstMetric(upload, prometheus.GaugeValue, math.Round(m_upload/1024/1024))

	log.Printf("Speedtest exporter finished")
}

func (e *Exporter) NetworkMetrics() (metrics Metrics, err error) {
	command := []string{"./speedtest-cli", "--json"}
	response, err := e.execute(strings.Join(command, " "))
	if err != nil {
		return metrics, err
	}

	err = json.Unmarshal(response.Output.Bytes(), &metrics)
	return metrics, err
}

// Response represents the command exec response.
type Response struct {
	Output bytes.Buffer
	Logs   bytes.Buffer
}

func (e *Exporter) execute(command string) (Response, error) {
	cmd := exec.Command("sh", "-c", fmt.Sprintf("%v", command))
	var outb, errb bytes.Buffer
	cmd.Stdout = &outb
	cmd.Stderr = &errb
	err := cmd.Run()
	if errb.Len() > 0 {
		// All logs are currently directed to stderr
		log.Printf("%v", errb.String())
	}
	return Response{Output: outb, Logs: errb}, err
}

func main() {
	var (
		addr = flag.String("listen-address", ":9112", "The address to listen on for HTTP requests.")
	)

	flag.Parse()

	log.Printf("Register exporter")
	prometheus.MustRegister(&Exporter{})

	log.Printf("Starting speedtest exporter")

	// Expose the registered metrics via HTTP.
	http.Handle("/metrics", promhttp.HandlerFor(
		prometheus.DefaultGatherer,
		promhttp.HandlerOpts{
			// Opt into OpenMetrics to support exemplars.
			EnableOpenMetrics: true,
		},
	))

	log.Printf("Listening on %v", *addr)

	// fatal
	log.Fatal(http.ListenAndServe(*addr, nil))
}
