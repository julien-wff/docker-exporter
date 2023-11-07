package exporter

import (
	"context"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"strings"
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
			Tag:        getImageTag(image.RepoTags, image.RepoDigests),
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

func getImageTag(repoTags []string, repoDigests []string) string {
	var tag string

	if len(repoTags) > 0 {
		tag = repoTags[0]
	} else if len(repoDigests) > 0 && strings.Contains(repoDigests[0], "@sha256:") {
		tag = strings.Split(repoDigests[0], "@sha256:")[0] + ":<none>"
	}

	if strings.TrimSpace(tag) == "" || tag == "<none>:<none>" {
		tag = "<none>"
	}

	return tag
}
