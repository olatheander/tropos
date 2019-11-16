package args

import (
	"github.com/docker/docker/pkg/homedir"
	"os"
	"path/filepath"
	"runtime"
)

type Docker struct {
	Image     string
	Workspace string
	Port      int32
}

type Kubernetes struct {
	Config         string
	Namespace      string
	DeploymentName string
	Image          string
	HostPort       int32
	ContainerName  string
	ContainerPort  int32
}

type Endpoint struct {
	Host string
	Port int
}

type SSH struct {
	User           string
	PrivateKeyPath string
	ServerEndpoint Endpoint
	LocalEndpoint  Endpoint
	RemoteEndpoint Endpoint
}

type Context struct {
	Docker     Docker
	Kubernetes Kubernetes
	SSH        SSH
}

//ParseArgs parse command line arguments
func ParseArgs() (Context, error) {
	/*
		if len(os.Args) == 2 {
			fmt.Println("expected 'stage' or 'proxy' sub-command")
			os.Exit(1)
		}

		stageCommand := flag.NewFlagSet("stage", flag.ExitOnError)
		namespace := stageCommand.String("namespace", "default", "namespace")
		sshPort := stageCommand.Int("ssh-port", 2022, "SSH port to expose")
		swapDeployment := stageCommand.String("swap-deployment", "", "Kubernetes deployment to swap out")

		proxyCommand := flag.NewFlagSet("proxy", flag.ExitOnError)

		switch os.Args[1] {
		case "stage":
			//User's CLI
			fmt.Println("subcommand 'proxy'")
			fmt.Println("  name:", *namespace)
			fmt.Println("  ssh-port:", *sshPort)
			fmt.Println("  Kubernetes deployment to swap out:", *swapDeployment)
			fmt.Println("  tail:", stageCommand.Args())

		case "proxy":
			// Running as proxy, i.e. typically inside the proxy container
			proxyCommand.Parse(os.Args[2:])
			fmt.Println("subcommand 'proxy'")
			fmt.Println("  tail:", proxyCommand.Args())
		default:
			fmt.Println("expected 'stage' or 'proxy' sub-commands")
			os.Exit(1)
		}
	*/
	var c Context
	c.Docker.Image = "docker.io/olatheander/tropos-base:latest"
	c.Docker.Workspace = "/tmp"                                           //TODO: get the path from cmd line params.
	c.Docker.Port = 2022                                                  //TODO: get the port from cmd line params.
	c.Kubernetes.Config = filepath.Join(homedir.Get(), ".kube", "config") //TODO: get from 1) cmd-line, 2) environment or 3) default.
	c.Kubernetes.DeploymentName = "tropos"                                //TODO: get the deployment name from 1) cmd-line, 2) environment or 3) use "tropos" as default.
	c.Kubernetes.Image = "docker.io/olatheander/tropos-base:latest"       //TODO: get the image from 1) cmd-line, 2) environment
	c.Kubernetes.HostPort = 8022                                          //TODO: get the port from 1) cmd-line, 2) environment or 3) 22 as default
	c.Kubernetes.ContainerPort = 22                                       //TODO: get the port from 1) cmd-line, 2) environment or 3) 22 as default
	c.SSH.User = "root"
	if runtime.GOOS == "windows" {
		//TODO: Consider https://github.com/mitchellh/go-homedir for this.
		c.SSH.PrivateKeyPath = filepath.Join(os.Getenv("HOME"), ".ssh", "id_rsa")
	} else {
		c.SSH.PrivateKeyPath = filepath.Join(os.Getenv("HOME"), ".ssh", "id_rsa")
	}
	c.SSH.ServerEndpoint.Host = "localhost"
	c.SSH.ServerEndpoint.Port = 8022
	c.SSH.LocalEndpoint.Host = "localhost"
	c.SSH.LocalEndpoint.Port = 2022
	c.SSH.RemoteEndpoint.Host = "localhost"
	c.SSH.RemoteEndpoint.Port = 10000
	return c, nil
}
