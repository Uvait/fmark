package cmd

import (
	"encoding/json"
	"fmt"
	"os"
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
	Use:   "add <name>",
	Short: "Add a bookmark",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		command, _ := cmd.Flags().GetString("command")
		text, _ := cmd.Flags().GetString("text")
		overwrite, _ := cmd.Flags().GetBool("overwrite")

		var value, path string
		if command != "" {
			value = command
			path, _ = xdg.DataFile("fmark/commands.json")
		} else {
			value = text
			path, _ = xdg.DataFile("fmark/texts.json")
		}

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

		newMark := MarkData{Name: args[0], Value: value}
		i := slices.IndexFunc(jsonData, func(m MarkData) bool {
			return m.Name == newMark.Name
		})

		if i != -1 {
			if overwrite {
				jsonData[i] = newMark
			} else {
				fmt.Fprintf(os.Stderr, "bookmark %q is already exists\n", newMark.Name)
				os.Exit(1)
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
	addCmd.Flags().StringP("command", "c", "", "save a command")
	addCmd.Flags().StringP("text", "t", "", "save a text")
	addCmd.Flags().BoolP("overwrite", "o", false, "save the bookmark, even if it already exists")
	addCmd.MarkFlagsOneRequired("command", "text")
	addCmd.MarkFlagsMutuallyExclusive("command", "text")
	rootCmd.AddCommand(addCmd)
}
