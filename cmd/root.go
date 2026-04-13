package cmd

import (
	"fmt"
	"os"

	"github.com/adrg/xdg"
	"github.com/spf13/afero"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:          "fmark",
	Short:        "Simple bookmark CLI",
	SilenceUsage: true,
}
var fs = afero.NewOsFs()

func Execute() {
	if xdg.DataHome == "" {
		fmt.Fprintln(os.Stderr, fmt.Errorf("Environment variable $XDG_DATA_HOME is empty"))
		os.Exit(1)
	} else if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
