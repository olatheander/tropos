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

	NewDeployment(&context.kubernetes)

	//TODO: Wait for SIGHUP and then clean up.
}
