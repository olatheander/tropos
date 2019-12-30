package main

import (
	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"os"
	"runtime"
	newDeploy "tropos/cmd/tropos/new"
	"tropos/cmd/tropos/swap"
)

func rootMain(command *cobra.Command, arguments []string) error {
	// If no commands were given, then print help information and bail.
	command.Help()

	// Success.
	return nil
}

var rootCommand = &cobra.Command{
	Use:          "tropos",
	Short:        "Tropos is a tool to enable native Kubernetes development with Visual Studio Code remote development features.",
	RunE:         rootMain,
	SilenceUsage: true,
}

var rootConfiguration struct {
	// help indicates whether or not help information should be shown for the
	// command.
	help bool
}

func init() {
	// Disable alphabetical sorting of commands in help output. This is a global
	// setting that affects all Cobra command instances.
	cobra.EnableCommandSorting = false

	// Disable Cobra's use of mousetrap. This breaks daemon registration on
	// Windows because it tries to enforce that the CLI only be launched from
	// a console, which it's not when running automatically.
	cobra.MousetrapHelpText = ""

	// Grab a handle for the command line flags.
	flags := rootCommand.Flags()

	// Disable alphabetical sorting of flags in help output.
	flags.SortFlags = false

	// Manually add a help flag to override the default message. Cobra will
	// still implement its logic automatically.
	flags.BoolVarP(&rootConfiguration.help, "help", "h", false, "Show help information")

	// Register commands.
	// HACK: Add the sync commands as direct subcommands of the root command for
	// temporary backward compatibility.
	commands := []*cobra.Command{
		newDeploy.NewCommand,
		swap.RootCommand,
		updateCommand,
	}
	rootCommand.AddCommand(commands...)

	// HACK If we're on Windows, enable color support for command usage and
	// error output by recursively replacing the output streams for Cobra
	// commands.
	if runtime.GOOS == "windows" {
		enableColorForCommand(rootCommand)
	}
}

// enableColorForCommand recursively enables colorized usage and error output
// for a command and all of its child commands.
func enableColorForCommand(command *cobra.Command) {
	// Enable color support for the command itself.
	command.SetOut(color.Output)
	command.SetErr(color.Error)

	// Recursively enable color support for child commands.
	for _, c := range command.Commands() {
		enableColorForCommand(c)
	}
}

func main() {
	// Execute the root command.
	if err := rootCommand.Execute(); err != nil {
		os.Exit(1)
	}
}
