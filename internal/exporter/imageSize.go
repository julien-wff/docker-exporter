package exporter

import (
	"context"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
)

type ImageSize struct {
	Id         string
	Tag        string
	Containers int
	Size       int
}

func ExportImageSize() []ImageSize {
	cli, err := client.NewClientWithOpts()
	if err != nil {
		panic(err)
	}

	images, err := cli.ImageList(context.Background(), types.ImageListOptions{})
	if err != nil {
		panic(err)
	}

	var imageSize []ImageSize
	for _, image := range images {
		imageSize = append(imageSize, ImageSize{
			Id:         image.ID,
			Tag:        image.RepoTags[0],
			Containers: 0,
			Size:       int(image.Size),
		})
	}

	containers, err := cli.ContainerList(context.Background(), types.ContainerListOptions{All: true})
	if err != nil {
		panic(err)
	}

	for _, container := range containers {
		for ind, image := range imageSize {
			if container.ImageID == image.Id {
				imageSize[ind].Containers++
				break
			}
		}
	}

	err = cli.Close()
	if err != nil {
		panic(err)
	}

	return imageSize
}
