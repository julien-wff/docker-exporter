package exporter

import (
	"context"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/volume"
	"github.com/docker/docker/client"
	"github.com/julien-wff/docker-exporter/internal/config"
	"log"
	"os/exec"
	"strconv"
	"strings"
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
		size, _ := getVolumeSize(vol.Mountpoint)
		volumeSize = append(volumeSize, VolumeSize{
			Name:           vol.Name,
			MountPoint:     vol.Mountpoint,
			Size:           size,
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

	err = cli.Close()
	if err != nil {
		panic(err)
	}

	return volumeSize
}

func getVolumeSize(mountpoint string) (int, error) {
	// Context with timeout
	cfg := config.GetConfig()
	ctx, cancel := context.WithTimeout(context.Background(), cfg.RequestTimeout)
	defer cancel()

	// Get the size of the volume using du
	cmd := exec.CommandContext(ctx, "du", "-sb", mountpoint)
	out, err := cmd.Output()
	if err != nil {
		log.Printf("error getting volume size at mountpoint %s: %s\n", mountpoint, err)
		return 0, err
	}

	// Parse the output of du
	size := string(out)
	size = strings.Split(size, "\t")[0]
	size = strings.Replace(size, "\n", "", -1)
	size = strings.Replace(size, " ", "", -1)

	// Return size converted to int
	intSize, err := strconv.Atoi(size)
	if err != nil {
		log.Printf("error converting volume size to int: %s\n", err)
		return 0, err
	}

	return intSize, nil
}
