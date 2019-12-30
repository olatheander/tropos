package args

import (
	"os"
	"path/filepath"
	"runtime"
)

const (
	DefaultDeployment    = "tropos"
	DefaultImage         = "docker.io/olatheander/tropos-base:latest"
	DefaultContainerPort = 22
	DefaultHost          = "localhost"
	DefaultHostPort      = 8022
	DefaultSshUser       = "root"
)

type Kubernetes struct {
	Context       string
	Namespace     string
	Deployment    string
	Image         string
	HostPort      uint16
	ContainerName string
	ContainerPort uint16
}

type SSH struct {
	User           string
	PrivateKeyPath string
}

type Context struct {
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
	//c.Kubernetes.Context = filepath.Join(homedir.Get(), ".kube", "config") //TODO: get from 1) cmd-line, 2) environment or 3) default.
	c.Kubernetes.Deployment = DefaultDeployment       //TODO: get the deployment name from 1) cmd-line, 2) environment or 3) use "tropos" as default.
	c.Kubernetes.Image = DefaultImage                 //TODO: get the image from 1) cmd-line, 2) environment
	c.Kubernetes.HostPort = DefaultHostPort           //TODO: get the port from 1) cmd-line, 2) environment or 3) 22 as default
	c.Kubernetes.ContainerPort = DefaultContainerPort //TODO: get the port from 1) cmd-line, 2) environment or 3) 22 as default
	c.SSH.User = DefaultSshUser
	if runtime.GOOS == "windows" {
		//TODO: Consider https://github.com/mitchellh/go-homedir for this.
		c.SSH.PrivateKeyPath = filepath.Join(os.Getenv("HOME"), ".ssh", "id_rsa")
	} else {
		c.SSH.PrivateKeyPath = filepath.Join(os.Getenv("HOME"), ".ssh", "id_rsa")
	}
	return c, nil
}
