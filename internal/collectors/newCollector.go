package collectors

import (
	"github.com/prometheus/client_golang/prometheus"
	"sync"
	"time"
)

type DockerCollector struct {
	volumesMetrics           *prometheus.Desc
	imageSizeMetrics         *prometheus.Desc
	networkContainersMetrics *prometheus.Desc
	scrapeDurationMetric     *prometheus.Desc
}

func (collector *DockerCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- collector.networkContainersMetrics
}

func (collector *DockerCollector) Collect(ch chan<- prometheus.Metric) {
	var wg sync.WaitGroup
	wg.Add(3)

	go CollectNetworkContainers(&wg, collector, ch)
	go CollectImageSize(&wg, collector, ch)
	go CollectVolumeSize(&wg, collector, ch)

	start := time.Now()
	wg.Wait()

	ch <- prometheus.MustNewConstMetric(collector.scrapeDurationMetric, prometheus.GaugeValue, time.Since(start).Seconds())
}

func NewDockerCollector() *DockerCollector {
	return &DockerCollector{
		volumesMetrics: prometheus.NewDesc(
			"docker_volume_size_bytes",
			"Size of docker volumes",
			[]string{"name", "mountpoint", "compose_project", "containers"},
			nil,
		),
		imageSizeMetrics: prometheus.NewDesc(
			"docker_image_size_bytes",
			"Size of docker images",
			[]string{"id", "tag", "containers"},
			nil,
		),
		networkContainersMetrics: prometheus.NewDesc(
			"docker_network_container_count",
			"Number of containers per network",
			[]string{"id", "name", "compose_project"},
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
