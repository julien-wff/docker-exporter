package exporter

import (
	"context"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/volume"
	"github.com/docker/docker/client"
	"log"
	"os"
	"path/filepath"
)

type VolumeSize struct {
	Name           string
	MountPoint     string
	Size           int
	Containers     int
	ComposeProject string
}

func ExportVolumeSize() []VolumeSize {
	cli, err := client.NewClientWithOpts()
	if err != nil {
		panic(err)
	}

	volumes, err := cli.VolumeList(context.Background(), volume.ListOptions{})
	if err != nil {
		panic(err)
	}

	var volumeSize []VolumeSize
	for _, vol := range volumes.Volumes {
		volumeSize = append(volumeSize, VolumeSize{
			Name:           vol.Name,
			MountPoint:     vol.Mountpoint,
			ComposeProject: vol.Labels["com.docker.compose.project"],
		})
	}

	containers, err := cli.ContainerList(context.Background(), types.ContainerListOptions{All: true})
	if err != nil {
		panic(err)
	}

	for _, container := range containers {
		for _, mount := range container.Mounts {
			for i, vol := range volumeSize {
				if vol.MountPoint == mount.Source {
					volumeSize[i].Containers++
				}
			}
		}
	}

	// Calculate up to 8 volumes in parallel
	sem := make(chan bool, 8)
	for i, vol := range volumeSize {
		sem <- true
		go func(i int, vol VolumeSize) {
			size, err := getVolumeSize(vol.MountPoint)
			if err != nil {
				log.Println("Error getting volume size:", err)
			}
			volumeSize[i].Size = size
			<-sem
		}(i, vol)
	}

	// Wait for all goroutines to finish
	for i := 0; i < cap(sem); i++ {
		sem <- true
	}

	err = cli.Close()
	if err != nil {
		panic(err)
	}

	return volumeSize
}

func getVolumeSize(mountpoint string) (int, error) {
	var size int
	err := filepath.Walk(mountpoint, func(_ string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			size += int(info.Size())
		}
		return nil
	})
	return size, err
}
