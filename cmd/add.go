package cmd

import (
	"encoding/json"
	"fmt"
	"slices"

	"github.com/adrg/xdg"
	"github.com/spf13/afero"
	"github.com/spf13/cobra"
)

type MarkData struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

var addCmd = &cobra.Command{
	Use:   "add <name> <command>",
	Short: "Add a bookmark",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		name := args[0]
		command := args[1]
		overwrite, _ := cmd.Flags().GetBool("overwrite")
		path, _ := xdg.DataFile("fmark/commands.json")

		if exists, _ := afero.Exists(fs, path); !exists {
			if err := afero.WriteFile(fs, path, []byte("[]"), 0644); err != nil {
				return err
			}
		}

		f, err := afero.ReadFile(fs, path)
		if err != nil {
			return err
		}

		var jsonData []MarkData
		if err := json.Unmarshal(f, &jsonData); err != nil {
			return err
		}

		newMark := MarkData{Name: name, Value: command}
		i := slices.IndexFunc(jsonData, func(m MarkData) bool {
			return m.Name == newMark.Name
		})
		if i != -1 {
			if overwrite {
				jsonData[i] = newMark
			} else {
				fmt.Printf("bookmark %q is already exists\n", newMark.Name)
				return nil
			}
		} else {
			jsonData = append(jsonData, newMark)
		}

		data, err := json.MarshalIndent(jsonData, "", "\t")
		if err != nil {
			return err
		}

		return afero.WriteFile(fs, path, data, 0644)
	},
}

func init() {
	addCmd.Flags().BoolP("overwrite", "o", false, "save the bookmark, even if it already exists")
	rootCmd.AddCommand(addCmd)
}
