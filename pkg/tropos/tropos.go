package tropos

import (
	"bytes"
	"fmt"
	"github.com/docker/docker/client"
	"io"
	"io/ioutil"
	appsv1 "k8s.io/api/apps/v1"
	"os"
	"os/exec"
	"os/signal"
	"strings"
	"sync"
	"tropos/pkg/args"
	"tropos/pkg/docker"
	"tropos/pkg/kubernetes"
	"tropos/pkg/ssh"
)

// Stage a new deployment and mount the workspace
func NewDeployment(context args.Context) {
	cli, err := docker.NewDockerClient()
	if err != nil {
		panic(err)
	}
	defer cli.Close()

	deployment, err := kubernetes.NewDeployment(&context.Kubernetes)
	if err != nil {
		panic(err)
	}
	defer func(deployment *appsv1.Deployment) {
		kubernetes.DeleteDeployment(&context.Kubernetes, deployment)
		fmt.Println("Deleted deployment (", deployment.Name, ")")
	}(deployment)

	err = authorizeSshKey(context.SSH.PublicKeyPath,
		&context.Kubernetes,
		deployment)
	if err != nil {
		panic(err)
	}

	err = generateSshKeysInPod(&context.Kubernetes, deployment)
	if err != nil {
		panic(err)
	}

	pubKeyFile, err := copyPodPublicKeyToTemp("/root/.ssh/tropos.pub",
		&context.Kubernetes,
		deployment)
	if err != nil {
		panic(err)
	}
	defer os.Remove(pubKeyFile.Name())

	fmt.Println("Starting Docker container")
	containerId, err := docker.CreateNewContainer(context.Docker.Image,
		context.Docker.Workspace,
		pubKeyFile.Name(),
		cli)
	if err != nil {
		panic(err)
	}
	fmt.Println("Started Docker container (ID:", containerId, ")")
	defer func(containerId string, cli *client.Client) {
		docker.CloseContainer(containerId, cli)
		fmt.Println("Stopped and removed Docker container (ID:", containerId, ")")
	}(containerId, cli)

	portForward(&context.Kubernetes, deployment, func(stopChannel chan struct{}) {
		defer close(stopChannel)

		fmt.Println("Mounting working directory in pod.")
		//TODO: this is just a dummy test. Should set up sshfs like in https://superuser.com/questions/616182/how-to-mount-local-directory-to-remote-like-sshfs
		var stdout, stderr bytes.Buffer
		err := kubernetes.Exec("ls -l /",
			&context.Kubernetes,
			deployment,
			nil,
			&stdout,
			&stderr)
		fmt.Println(&stdout)
		fmt.Println(&stderr)
		if err != nil {
			panic(err)
		}

		fmt.Println("Configuring SSH tunnels.")
		setupSsh(&context.SSH)

		fmt.Println("All set up. Carry on... (press Ctrl+C to exit).")
		waitForCtrlC()
		fmt.Println("Closing down and cleaning up.")
	})
}

func setupSsh(sshConfig *args.SSH) {
	ssh.NewSSHTunnel(sshConfig.User,
		sshConfig.PublicKeyPath,
		&ssh.Endpoint{
			Host: sshConfig.ServerEndpoint.Host,
			Port: sshConfig.ServerEndpoint.Port,
		},
		&ssh.Endpoint{
			Host: sshConfig.LocalEndpoint.Host,
			Port: sshConfig.LocalEndpoint.Port,
		},
		&ssh.Endpoint{
			Host: sshConfig.RemoteEndpoint.Host,
			Port: sshConfig.RemoteEndpoint.Port,
		})
}

func waitForCtrlC() {
	var endWaiter sync.WaitGroup
	endWaiter.Add(1)
	var signalChannel chan os.Signal
	signalChannel = make(chan os.Signal, 1)
	signal.Notify(signalChannel, os.Interrupt)
	go func() {
		<-signalChannel
		endWaiter.Done()
	}()
	endWaiter.Wait()
}

// Swap out an existing deployment for a new development deployment
func SwapDeployment(context args.Context) {

}

// Authorize the SSH key in the deployment.
// Copy SSH public key to container, i.e. the equivalent of kubectl cp ~/.ssh/id_rsa.pub tropos-58d96c958d-d4799:/root/.ssh/authorized_keys
func authorizeSshKey(keyPath string, k8s *args.Kubernetes, deployment *appsv1.Deployment) (error) {
	reader, writer := io.Pipe()

	defer writer.Close()
	cmd := exec.Command("cat", keyPath)
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
		if err != nil {
			panic(err)
		}

		// Change ovnership of the authorized_keys file.
		err = kubernetes.Exec("chown root:root /root/.ssh/authorized_keys",
			k8s,
			deployment,
			reader,
			&stdout,
			&stderr)
		if err != nil {
			panic(err)
		}

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

// Make the Docker container mounting the workspace files trust the new SSH key in the Pod
func containerTrustPodKeys(k8s *args.Kubernetes,
	deployment *appsv1.Deployment,
	containerId string,
	cli *client.Client) {

	file, err := ioutil.TempFile("", "tropos.*.pub")
	if err != nil {
		panic(err)
	}
	defer os.Remove(file.Name())

	copyPodPublicKey("/root/.ssh/tropos.pub",
		file.Name(),
		k8s,
		deployment)

	docker.CopyToContainer(containerId,
		file.Name(),
		"/root/.ssh/authorized_keys",
		cli)
	fmt.Println("Authorized SSH key in Docker container:")
}

// Copy the public SSH key from the Pod to a temporary file
func copyPodPublicKeyToTemp(keyPath string,
	k8s *args.Kubernetes,
	deployment *appsv1.Deployment) (f *os.File, err error) {

	file, err := ioutil.TempFile("", "tropos.*.pub")
	if err != nil {
		panic(err)
	}

	err = copyPodPublicKey(keyPath,
		file.Name(),
		k8s,
		deployment)

	return file, err
}

// Copy the public SSH key from the Pod
func copyPodPublicKey(keyPath string,
	destPath string,
	k8s *args.Kubernetes,
	deployment *appsv1.Deployment) error {
	reader, writer := io.Pipe()

	defer reader.Close()
	var stdout, stderr bytes.Buffer
	cmd := exec.Command("tee", destPath)
	cmd.Stdin = reader
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	go func() {
		defer writer.Close()
		var stderr bytes.Buffer
		cmd := []string{"cat", keyPath}
		err := kubernetes.Exec(strings.Join(cmd, " "),
			k8s,
			deployment,
			nil,
			writer,
			&stderr)

		if err != nil {
			panic(err)
		}
	}()

	cmd.Run()
	fmt.Println("Copied SSH key from Pod:", &stdout)
	return nil
}

// Set up port-forwarding and trigger f-function after port-forwarding is setup.
func portForward(k8s *args.Kubernetes, deployment *appsv1.Deployment, f func(stopChannel chan struct{})) error {
	stopChan, readyChan := make(chan struct{}, 1), make(chan struct{}, 1)
	out, errOut := new(bytes.Buffer), new(bytes.Buffer)
	go func() {
		for range readyChan { // Kubernetes will close this channel when it has something to tell us.
		}
		if len(errOut.String()) != 0 {
			panic(errOut.String())
		} else if len(out.String()) != 0 {
			fmt.Println(out.String())
			f(stopChan)
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
