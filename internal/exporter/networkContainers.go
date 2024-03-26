package exporter

import (
	"context"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
)

type NetworkContainers struct {
	Id             string
	Network        string
	Containers     int
	ComposeProject string
}

func ExportNetworkContainers() []NetworkContainers {
	cli, err := client.NewClientWithOpts()
	if err != nil {
		panic(err)
	}

	networks, err := cli.NetworkList(context.Background(), types.NetworkListOptions{})
	if err != nil {
		panic(err)
	}

	var networkContainers []NetworkContainers
	for _, network := range networks {
		networkContainers = append(networkContainers, NetworkContainers{
			Id:             network.ID,
			Network:        network.Name,
			Containers:     0,
			ComposeProject: network.Labels["com.docker.compose.project"],
		})
	}

	containers, err := cli.ContainerList(context.Background(), container.ListOptions{All: true})
	if err != nil {
		panic(err)
	}

	for _, ctr := range containers {
		for _, containerNetwork := range ctr.NetworkSettings.Networks {
			for netInd, networkContainer := range networkContainers {
				if networkContainer.Id == containerNetwork.NetworkID {
					networkContainers[netInd].Containers++
				}
			}
		}
	}

	err = cli.Close()
	if err != nil {
		panic(err)
	}

	return networkContainers
}
