package tropos

import (
	"bytes"
	"fmt"
	"io"
	appsv1 "k8s.io/api/apps/v1"
	"os/exec"
	"tropos/pkg/args"
	"tropos/pkg/docker"
	"tropos/pkg/kubernetes"
)

// Stage a new deployment and mount the workspace
func NewDeployment(context args.Context) {
	cli, err := docker.NewDockerClient()
	if err != nil {
		panic(err)
	}

	fmt.Println("Starting Docker container")
	containerId, err := docker.CreateNewContainer(context.Docker.Image,
		context.Docker.Workspace,
		cli)
	if err != nil {
		panic(err)
	}
	fmt.Println("Started Docker container (ID:", containerId, ")")

	deployment, err := kubernetes.NewDeployment(&context.Kubernetes)
	if err != nil {
		panic(err)
	}

	authorizeSshKey(&context.Kubernetes, deployment)

	generateSshKeysInPod(&context.Kubernetes, deployment)

	portForward(&context.Kubernetes, deployment)

	fmt.Println("All set up. Carry on...")
	//TODO: Wait for SIGHUP and then clean up.
}

// Swap out an existing deployment for a new development deployment
func SwapDeployment(context args.Context) {

}

// Authorize the SSH key in the deployment.
// Copy SSH public key to container, i.e. the equivalent of kubectl cp ~/.ssh/id_rsa.pub tropos-58d96c958d-d4799:/root/.ssh/authorized_keys
func authorizeSshKey(k8s *args.Kubernetes, deployment *appsv1.Deployment) (error) {
	reader, writer := io.Pipe()

	defer writer.Close()
	cmd := exec.Command("cat", "/home/olathe/.ssh/id_rsa.pub")
	cmd.Stdout = writer

	go func() {
		defer reader.Close()
		var stdout, stderr bytes.Buffer
		err := kubernetes.Exec("tee /root/.ssh/authorized_keys",
			k8s,
			deployment,
			reader,
			&stdout,
			&stderr)

		fmt.Println("Authorized SSH key:", &stdout)

		if err != nil {
			panic(err)
		}
	}()

	cmd.Run()
	return nil
}

// Generate a new SSH key pair in the Pod using ssh-keygen.
func generateSshKeysInPod(k8s *args.Kubernetes, deployment *appsv1.Deployment) (error) {

	var stdout, stderr bytes.Buffer
	err := kubernetes.Exec("/scripts/generate-ssh-keys.sh",
		k8s,
		deployment,
		nil,
		&stdout,
		&stderr)

	if err != nil {
		fmt.Println(&stdout)
		fmt.Println(&stderr)
		panic(err)
	}

	fmt.Println("Generated Tropos SSH keys in Pod")
	return nil
}

// Set up port-forwarding.
func portForward(k8s *args.Kubernetes, deployment *appsv1.Deployment) (error) {
	stopChan, readyChan := make(chan struct{}, 1), make(chan struct{}, 1)
	out, errOut := new(bytes.Buffer), new(bytes.Buffer)
	go func() {
		for range readyChan { // Kubernetes will close this channel when it has something to tell us.
		}
		if len(errOut.String()) != 0 {
			panic(errOut.String())
		} else if len(out.String()) != 0 {
			fmt.Println(out.String())
			go func() {
				fmt.Println("Mounting working directory in pod.")
				//TODO: this is just a dummy test. Should set up sshfs like in https://superuser.com/questions/616182/how-to-mount-local-directory-to-remote-like-sshfs
				var stdout, stderr bytes.Buffer
				err := kubernetes.Exec("ls -l /",
					k8s,
					deployment,
					nil,
					&stdout,
					&stderr)
				fmt.Println(&stdout)
				fmt.Println(&stderr)
				if err != nil {
					panic(err)
				}
			}()
		}
	}()

	err := kubernetes.PortForward(k8s,
		deployment,
		readyChan,
		stopChan,
		out,
		errOut)
	if err != nil {
		panic(err)
	}

	return nil
}
