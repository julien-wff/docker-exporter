package exporter

import (
	"context"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/image"
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

	images, err := cli.ImageList(context.Background(), image.ListOptions{})
	if err != nil {
		panic(err)
	}

	var imageSize []ImageSize
	for _, img := range images {
		imageSize = append(imageSize, ImageSize{
			Id:         img.ID,
			Tag:        getImageTag(img.RepoTags, img.RepoDigests),
			Containers: 0,
			Size:       int(img.Size),
		})
	}

	containers, err := cli.ContainerList(context.Background(), container.ListOptions{All: true})
	if err != nil {
		panic(err)
	}

	for _, ctr := range containers {
		for ind, img := range imageSize {
			if ctr.ImageID == img.Id {
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
