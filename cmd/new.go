/*
Copyright © 2020 NAME HERE <EMAIL ADDRESS>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"github.com/docker/docker/pkg/homedir"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"os"
	"path/filepath"
	"runtime"
	"tropos/pkg/args"
	"tropos/pkg/tropos"
)

// newCmd represents the new command
var newCmd = &cobra.Command{
	Use:   "new",
	Short: "Create a new Tropos enabled deployment.",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	RunE: newMain,
}

var newConfiguration struct {
	Context args.Context
}

func init() {
	rootCmd.AddCommand(newCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// newCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// newCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	flags := newCmd.Flags()
	flags.StringVarP(&newConfiguration.Context.Kubernetes.Deployment,
		"deployment",
		"d",
		args.DefaultDeployment,
		"Specify the name of the Kubernetes deployment")
	flags.StringVarP(&newConfiguration.Context.Kubernetes.Namespace,
		"namespace",
		"n",
		"",
		"Specify the Kubernetes namespace")
	flags.StringVarP(&newConfiguration.Context.Kubernetes.Image,
		"image",
		"",
		args.DefaultImage,
		"Specify the Docker image to wire up in the deployment")
	flags.Int32Var(&newConfiguration.Context.Kubernetes.ContainerPort,
		"container-port",
		args.DefaultContainerPort,
		"Specify the SSH port of the deployment image")
	flags.Int32Var(&newConfiguration.Context.Kubernetes.HostPort,
		"host-port",
		args.DefaultHostPort,
		"Specify the mapped SSH port on the host")
	flags.StringVarP(&newConfiguration.Context.Kubernetes.Context,
		"context",
		"c",
		"",
		"Specify the Kubernetes context")
	flags.StringVar(&newConfiguration.Context.Kubernetes.Config,
		"config",
		getDefaultK8sConfigPath(),
		"Specify the Kubernetes configuration file")
	flags.StringVarP(&newConfiguration.Context.SSH.User,
		"ssh-user",
		"u",
		args.DefaultSshUser,
		"Specify the SSH user for connecting to the deployment")
	flags.StringVarP(&newConfiguration.Context.SSH.IdentityFile,
		"identity-file",
		"i",
		getDefaultSshKeyPath(),
		"Specify the SSH key to use for authentication")
}

func getDefaultSshKeyPath() string {
	if runtime.GOOS == "windows" {
		//TODO: Consider https://github.com/mitchellh/go-homedir for this.
		return filepath.Join(os.Getenv("HOME"), ".ssh", "id_rsa")
	} else {
		return filepath.Join(os.Getenv("HOME"), ".ssh", "id_rsa")
	}
}

func getDefaultK8sConfigPath() string {
	//TODO: get from 1) cmd-line, 2) environment or 3) default.
	return filepath.Join(homedir.Get(), ".kube", "config")
}

func newMain(command *cobra.Command, arguments []string) error {
	// Validate arguments.
	if len(arguments) != 0 {
		return errors.New("unexpected arguments provided")
	}

	//flags := command.Flags()
	//
	//deployment, err := flags.GetString("deployment")
	//if err != nil {
	//	return err
	//}
	//newConfiguration.Context.Kubernetes.Deployment = deployment
	//
	//ns, err := flags.GetString("namespace")
	//if err != nil {
	//	return err
	//}
	//newConfiguration.Context.Kubernetes.Namespace = ns
	//
	//image, err := flags.GetString("image")
	//if err != nil {
	//	return err
	//}
	//newConfiguration.Context.Kubernetes.Image = image
	//
	//containerPort, err := flags.GetInt32("container-port")
	//if err != nil {
	//	return err
	//}
	//newConfiguration.Context.Kubernetes.ContainerPort = containerPort
	//
	//hostPort, err := flags.GetInt32("host-port")
	//if err != nil {
	//	return err
	//}
	//newConfiguration.Context.Kubernetes.HostPort = hostPort
	//
	//context, err := flags.GetString("context")
	//if err != nil {
	//	return err
	//}
	//newConfiguration.Context.Kubernetes.Context = context
	//
	//config, err := flags.GetString("config")
	//if err != nil {
	//	return err
	//}
	//newConfiguration.Context.Kubernetes.Config = config
	//
	//sshUser, err := flags.GetString("ssh-user")
	//if err != nil {
	//	return err
	//}
	//newConfiguration.Context.SSH.User = sshUser
	//
	//identityFile, err := flags.GetString("identity-file")
	//if err != nil {
	//	return err
	//}
	//newConfiguration.Context.SSH.IdentityFile = identityFile

	return tropos.NewDeployment(newConfiguration.Context)
}
