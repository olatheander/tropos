package docker

import (
	"context"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/client"
	"github.com/docker/docker/pkg/archive"
	"github.com/docker/docker/pkg/jsonmessage"
	"github.com/docker/docker/pkg/term"
	"github.com/docker/go-connections/nat"
	"os"
)

func NewDockerClient() (*client.Client, error) {
	return client.NewClientWithOpts(client.WithVersion("1.39")) //TODO: Docker API version as cmd line option or negotiate with server?
}

func CreateNewContainer(image string,
	workspaceLocalPath string,
	pubKeyLocalPath string,
	cli *client.Client) (string, error) {
	hostBinding := nat.PortBinding{
		HostIP:   "0.0.0.0",
		HostPort: "2022",
	}
	containerPort, err := nat.NewPort("tcp", "22")
	if err != nil {
		panic(err)
	}

	portBinding := nat.PortMap{containerPort: []nat.PortBinding{hostBinding}}
	PullImage(image, cli)
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
				{
					Type:   mount.TypeBind,
					Source: pubKeyLocalPath,
					Target: "/root/.ssh/authorized_keys",
				},
			},
		},
		nil,
		"")
	if err != nil {
		panic(err)
	}

	cli.ContainerStart(context.Background(), cont.ID, types.ContainerStartOptions{})
	cli.ContainerWait(context.Background(), cont.ID, container.WaitConditionNotRunning)
	//statusCh, errCh := cli.ContainerWait(context.Background(), cont.ID, container.WaitConditionNotRunning)
	//select {
	//case err := <-errCh:
	//	if err != nil {
	//		panic(err)
	//	}
	//case <-statusCh:
	//}
	return cont.ID, nil
}

func PullImage(image string, cli *client.Client) error {
	reader, err := cli.ImagePull(context.Background(),
		image,
		types.ImagePullOptions{
			All: false,
		})
	if err != nil {
		panic(err)
	}
	defer reader.Close()
	termFd, isTerm := term.GetFdInfo(os.Stderr)
	jsonmessage.DisplayJSONMessagesStream(reader, os.Stderr, termFd, isTerm, nil)
	return nil
}

func CloseContainer(containerId string,
	cli *client.Client) error {
	err := cli.ContainerStop(context.Background(),
		containerId,
		nil)
	if err != nil {
		panic(err)
	}
	err = cli.ContainerRemove(context.Background(),
		containerId,
		types.ContainerRemoveOptions{
			RemoveVolumes: false,
			RemoveLinks:   false,
			Force:         false,
		})
	if err != nil {
		panic(err)
	}
	return nil
}

func CopyToContainer(containerId string,
	srcPath string,
	destPath string,
	cli *client.Client) error {

	// Prepare source copy info.
	srcInfo, err := archive.CopyInfoSourcePath(srcPath, true)
	if err != nil {
		return err
	}

	srcArchive, err := archive.TarResource(srcInfo)
	if err != nil {
		return err
	}
	defer srcArchive.Close()

	err = cli.CopyToContainer(context.Background(),
		containerId,
		destPath,
		srcArchive,
		types.CopyToContainerOptions{AllowOverwriteDirWithFile: false,
			CopyUIDGID: false})
	if err != nil {
		panic(err)
	}

	return nil
}
