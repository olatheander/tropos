package main

import (
	"tropos/pkg/args"
	"tropos/pkg/tropos"
)

func main() {
	context, err := args.ParseArgs()
	if err != nil {
		panic(err)
	}

	tropos.NewDeployment(context)
}
