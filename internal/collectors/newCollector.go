package collectors

import (
	"github.com/prometheus/client_golang/prometheus"
	"sync"
	"time"
)

type DockerCollector struct {
	networkContainersMetrics *prometheus.Desc
	scrapeDurationMetric     *prometheus.Desc
}

func (collector *DockerCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- collector.networkContainersMetrics
}

func (collector *DockerCollector) Collect(ch chan<- prometheus.Metric) {
	var wg sync.WaitGroup
	wg.Add(1)

	go CollectNetworkContainers(&wg, collector, ch)

	start := time.Now()
	wg.Wait()

	ch <- prometheus.MustNewConstMetric(collector.scrapeDurationMetric, prometheus.GaugeValue, time.Since(start).Seconds())
}

func NewDockerCollector() *DockerCollector {
	return &DockerCollector{
		networkContainersMetrics: prometheus.NewDesc(
			"docker_network_container_count",
			"Number of containers per network",
			[]string{"name"},
			nil,
		),
		scrapeDurationMetric: prometheus.NewDesc(
			"docker_last_scrape_duration_seconds",
			"Time it took to scrape docker metrics",
			nil,
			nil,
		),
	}
}
