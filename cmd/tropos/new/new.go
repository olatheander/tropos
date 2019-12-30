package new

import (
	"github.com/spf13/cobra"
	"tropos/pkg/args"
)

func newMain(command *cobra.Command, arguments []string) error {
	// If no commands were given, then print help information and bail. We don't
	// have to worry about warning about arguments being present here (which
	// would be incorrect usage) because arguments can't even reach this point
	// (they will be mistaken for subcommands and a error will be displayed).
	command.Help()

	// Success.
	return nil
}

var NewCommand = &cobra.Command{
	Use:          "new",
	Short:        "Create a new Tropos enabled deployment.",
	RunE:         newMain,
	SilenceUsage: true,
}

var newConfiguration struct {
	// help indicates whether or not help information should be shown for the
	// command.
	help    bool
	Context args.Context
}

func init() {
	// Grab a handle for the command line flags.
	flags := NewCommand.Flags()

	// Disable alphabetical sorting of flags in help output.
	flags.SortFlags = false

	// Manually add a help flag to override the default message. Cobra will
	// still implement its logic automatically.
	flags.BoolVarP(&newConfiguration.help, "help", "h", false, "Show help information")

	flags.StringVarP(&newConfiguration.Context.Kubernetes.Context,
		"deployment",
		"d",
		args.DefaultDeployment,
		"Specify the name of the Kubernetes deployment")
	flags.StringVarP(&newConfiguration.Context.Kubernetes.Namespace,
		"namespace",
		"n",
		"default",
		"Specify the Kubernetes namespace")
	flags.StringVarP(&newConfiguration.Context.Kubernetes.Image,
		"image",
		"",
		args.DefaultImage,
		"Specify the Docker image to wire up in the deployment")
	flags.Uint16Var(&newConfiguration.Context.Kubernetes.ContainerPort,
		"container-port",
		args.DefaultContainerPort,
		"Specify the SSH port of the deployment image")
	flags.Uint16Var(&newConfiguration.Context.Kubernetes.HostPort,
		"host-port",
		args.DefaultHostPort,
		"Specify the mapped SSH port on the host")
	flags.StringVarP(&newConfiguration.Context.SSH.User,
		"ssh-user",
		"u",
		args.DefaultSshUser,
		"Specify the SSH user for connecting to the deployment")
	flags.StringVarP(&newConfiguration.Context.SSH.PrivateKeyPath,
		"identity-file",
		"i",
		"~/.ssh/id_rsa",
		"Specify the SSH key to use for authentication")
	flags.StringVarP(&newConfiguration.Context.Kubernetes.Context,
		"context",
		"c",
		"",
		"Specify the Kubernetes context")
}
