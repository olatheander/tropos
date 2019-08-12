package main

func main() {
	context, err := ParseArgs()
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

	deployment, err := NewDeployment(&context.kubernetes)
	if err != nil {
		panic(err)
	}

	//deployment, err := GetDeployment(&context.kubernetes)
	//_, err = GetDeploymentPods(&context.kubernetes, deployment)
	err = PortForward(&context.kubernetes,
		deployment)
	if err != nil {
		panic(err)
	}

	//TODO: Wait for SIGHUP and then clean up.
}
