package cmd

import (
	"fmt"
	"slices"

	"github.com/spf13/cobra"
)

var showCmd = &cobra.Command{
	Use:               "show [<name>...]",
	Short:             "Print bookmarks",
	ValidArgsFunction: getBookmarkValidArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		all, _ := cmd.Flags().GetBool("all")
		if !all && len(args) == 0 {
			return fmt.Errorf("specify bookmark name or use --all")
		}
		if jsonData, err := loadMarks(); err != nil {
			return err
		} else {
			if all {
				if len(jsonData) == 0 {
					return fmt.Errorf("you dont have any bookmark")
				}
				for _, bookmark := range jsonData {
					fmt.Println(bookmark)
				}
			} else {
				for _, name := range args {
					if i := slices.IndexFunc(jsonData, func(m MarkData) bool {
						return name == m.Name
					}); i != -1 {
						fmt.Println(jsonData[i])
					} else {
						return fmt.Errorf("bookmark %q not found", name)
					}
				}
			}
			return nil
		}
	},
}

func init() {
	showCmd.Flags().BoolP("all", "a", false, "print all bookmarks")
	rootCmd.AddCommand(showCmd)
}
