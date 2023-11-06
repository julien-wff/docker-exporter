package exporter

import (
	"context"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
)

type NetworkContainers struct {
	Network    string
	Containers int
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
			Network: network.Name,
			// TODO: the containers field is always empty, need to find a way to get the containers in a network
			Containers: len(network.Containers),
		})
	}

	err = cli.Close()
	if err != nil {
		panic(err)
	}

	return networkContainers
}
