package collectors

import (
	"github.com/julien-wff/docker-exporter/internal/exporter"
	"github.com/prometheus/client_golang/prometheus"
	"strconv"
	"sync"
)

func CollectImageSize(wg *sync.WaitGroup, collector *DockerCollector, ch chan<- prometheus.Metric) {
	defer wg.Done()

	imageSize := exporter.ExportImageSize()

	for _, img := range imageSize {
		ch <- prometheus.MustNewConstMetric(
			collector.imageSizeMetrics,
			prometheus.GaugeValue,
			float64(img.Size),
			img.Id,
			img.Tag,
			strconv.Itoa(img.Containers))
	}
}
