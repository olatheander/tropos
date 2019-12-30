package main

import (
	"fmt"

	"github.com/spf13/cobra"
)

func versionMain(command *cobra.Command, arguments []string) error {
	// Print version information.
	fmt.Println("0.0.1")	//TODO: show proper version.

	// Success.
	return nil
}

var versionCommand = &cobra.Command{
	Use:          "version",
	Short:        "Show version information",
	RunE:         versionMain,
	SilenceUsage: true,
}

var versionConfiguration struct {
	// help indicates whether or not help information should be shown for the
	// command.
	help bool
}

func init() {
	// Grab a handle for the command line flags.
	flags := versionCommand.Flags()

	// Disable alphabetical sorting of flags in help output.
	flags.SortFlags = false

	// Manually add a help flag to override the default message. Cobra will
	// still implement its logic automatically.
	flags.BoolVarP(&versionConfiguration.help, "help", "h", false, "Show help information")
}
