package docker

import (
	"context"
	"fmt"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"
)

func NewDockerClient() (*client.Client, error) {
	return client.NewClientWithOpts(client.WithVersion("1.39")) //TODO: Docker API version as cmd line option or negotiate with server?
}

func CreateNewContainer(image string,
	workspaceLocalPath string,
	cli *client.Client) (string, error) {
	fmt.Println("Starting a new container.")
	hostBinding := nat.PortBinding{
		HostIP:   "0.0.0.0",
		HostPort: "2022",
	}
	containerPort, err := nat.NewPort("tcp", "22")
	if err != nil {
		panic("Unable to get the port")
	}

	portBinding := nat.PortMap{containerPort: []nat.PortBinding{hostBinding}}
	cont, err := cli.ContainerCreate(
		context.Background(),
		&container.Config{
			Image: image,
		},
		&container.HostConfig{
			PortBindings: portBinding,
			Mounts: []mount.Mount{
				{
					Type:   mount.TypeBind,
					Source: workspaceLocalPath,
					Target: "/workspace",
				},
			},
		},
		nil,
		"")
	if err != nil {
		panic(err)
	}

	cli.ContainerStart(context.Background(), cont.ID, types.ContainerStartOptions{})
	fmt.Printf("Container %s is started", cont.ID)
	return cont.ID, nil
}
