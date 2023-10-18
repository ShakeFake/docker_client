package utils

import (
	"context"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"io"
)

type DockerImpl interface {
	ImagePull(imageName string) (io.ReadCloser, error)
	ContainerCreate(config *container.Config) (container.ContainerCreateCreatedBody, error)
}

type DockerClient struct {
	Client *client.Client
	// Ctx 先释放出去，暂时不写对应函数。
	Ctx context.Context
}

var CLIENT DockerClient

func init() {
	ctx := context.Background()
	var err error

	CLIENT.Client, err = client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		panic(err)
	}
	CLIENT.Ctx = ctx
}

func (d *DockerClient) ImagePull(imageName string) (io.ReadCloser, error) {
	reader, err := d.Client.ImagePull(d.Ctx, imageName, types.ImagePullOptions{})
	return reader, err
}

func (d *DockerClient) ContainerCreate(config *container.Config) (container.ContainerCreateCreatedBody, error) {
	resp, err := d.Client.ContainerCreate(d.Ctx, config, nil, nil, nil, "")
	return resp, err
}
