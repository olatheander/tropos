package main

import (
	"github.com/docker/docker/pkg/homedir"
	"path/filepath"
)

type Proxy struct {
	image     string
	workspace string
	port      int32
}

type Kubernetes struct {
	config         string
	namespace      string
	deploymentName string
	image          string
	hostPort       int32
	containerName  string
	containerPort  int32
}

type Context struct {
	proxy      Proxy
	kubernetes Kubernetes
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
	c.proxy.workspace = "/tmp"                                            //TODO: get the path from cmd line params.
	c.proxy.port = 2022                                                   //TODO: get the port from cmd line params.
	c.kubernetes.config = filepath.Join(homedir.Get(), ".kube", "config") //TODO: get from 1) cmd-line, 2) environment or 3) default.
	c.kubernetes.deploymentName = "tropos"                                //TODO: get the deployment name from 1) cmd-line, 2) environment or 3) use "tropos" as default.
	//c.kubernetes.image = "nginx:1.12"                                   //TODO: get the image from 1) cmd-line, 2) environment
	c.kubernetes.image = "tropos-base:latest"                            //TODO: get the image from 1) cmd-line, 2) environment
	c.kubernetes.hostPort = 2022                                          //TODO: get the port from 1) cmd-line, 2) environment or 3) 22 as default
	c.kubernetes.containerPort = 22                                       //TODO: get the port from 1) cmd-line, 2) environment or 3) 22 as default

	return c, nil
}
