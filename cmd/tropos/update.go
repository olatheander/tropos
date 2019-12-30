package main

import (
	"fmt"

	"github.com/spf13/cobra"
)

func updateMain(command *cobra.Command, arguments []string) error {
	// Print version information.
	fmt.Println("Troops is already the latest version")

	// Success.
	return nil
}

var updateCommand = &cobra.Command{
	Use:          "update",
	Short:        "Update Tropos to the latest version",
	RunE:         updateMain,
	SilenceUsage: true,
}

var updateConfiguration struct {
	// help indicates whether or not help information should be shown for the
	// command.
	help bool
}

func init() {
	// Grab a handle for the command line flags.
	flags := updateCommand.Flags()

	// Disable alphabetical sorting of flags in help output.
	flags.SortFlags = false

	// Manually add a help flag to override the default message. Cobra will
	// still implement its logic automatically.
	flags.BoolVarP(&updateConfiguration.help, "help", "h", false, "Show help information")
}
