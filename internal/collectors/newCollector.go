package collectors

import (
	"github.com/prometheus/client_golang/prometheus"
	"sync"
)

type DockerCollector struct {
	networkContainersMetrics *prometheus.Desc
}

func (collector *DockerCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- collector.networkContainersMetrics
}

func (collector *DockerCollector) Collect(ch chan<- prometheus.Metric) {
	var wg sync.WaitGroup
	wg.Add(1)
	go CollectNetworkContainers(&wg, collector, ch)
	wg.Wait()
}

func NewDockerCollector() *DockerCollector {
	return &DockerCollector{
		networkContainersMetrics: prometheus.NewDesc(
			"docker_network_container_count",
			"Number of containers per network",
			[]string{"name"},
			nil,
		),
	}
}
