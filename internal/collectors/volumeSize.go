package collectors

import (
	"github.com/julien-wff/docker-exporter/internal/exporter"
	"github.com/prometheus/client_golang/prometheus"
	"strconv"
	"sync"
)

func CollectVolumeSize(wg *sync.WaitGroup, collector *DockerCollector, ch chan<- prometheus.Metric) {
	defer wg.Done()

	volumesSize := exporter.ExportVolumeSize()

	for _, vol := range volumesSize {
		ch <- prometheus.MustNewConstMetric(
			collector.volumesMetrics,
			prometheus.GaugeValue,
			float64(vol.Size),
			vol.Name,
			vol.MountPoint,
			strconv.Itoa(vol.Containers))
	}
}
