package collectors

import (
	"github.com/julien-wff/docker-exporter/internal/exporter"
	"github.com/prometheus/client_golang/prometheus"
	"sync"
)

func CollectNetworkContainers(wg *sync.WaitGroup, collector *DockerCollector, ch chan<- prometheus.Metric) {
	defer wg.Done()

	networkContainers := exporter.ExportNetworkContainers()

	for _, net := range networkContainers {
		ch <- prometheus.MustNewConstMetric(
			collector.networkContainersMetrics,
			prometheus.GaugeValue,
			float64(net.Containers),
			net.Network)
	}
}
