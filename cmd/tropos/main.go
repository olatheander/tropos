package main

import (
	"tropos/pkg/args"
	"tropos/pkg/kubernetes"
)

func main() {
	context, err := args.ParseArgs()
	if err != nil {
		panic(err)
	}

	//cli, err := NewDockerClient()
	//if err != nil {
	//	panic(err)
	//}

	//CreateNewContainer("tropos-proxy",
	//	context.proxy.workspace,
	//	cli)

	deployment, err := kubernetes.NewDeployment(&context.Kubernetes)
	if err != nil {
		panic(err)
	}

	//deployment, err := GetDeployment(&context.kubernetes)
	//_, err = GetDeploymentPods(&context.kubernetes, deployment)
	err = kubernetes.PortForward(&context.Kubernetes,
		deployment)
	if err != nil {
		panic(err)
	}

	//TODO: Wait for SIGHUP and then clean up.
}
