package cmd

import (
	"encoding/json"
	"slices"

	"github.com/adrg/xdg"
	"github.com/spf13/afero"
	"github.com/spf13/cobra"
)

var deleteCmd = &cobra.Command{
	Use:     "delete <name>...",
	Aliases: []string{"del"},
	Short:   "Delete a bookmark",
	Args:    cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		path, _ := xdg.DataFile("fmark/commands.json")
		if jsonData, err := loadMarks(); err != nil {
			return err
		} else {
			for _, name := range args {
				jsonData = slices.DeleteFunc(jsonData, func(m MarkData) bool {
					return name == m.Name
				})
			}
			if data, err := json.MarshalIndent(jsonData, "", "\t"); err != nil {
				return err
			} else {
				return afero.WriteFile(fs, path, data, 0644)
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(deleteCmd)
}
