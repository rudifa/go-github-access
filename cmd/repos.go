/*
Copyright © 2023 @rudifa Rudolf Farkas rudi.farkas@gmail.com
License MIT
*/
package cmd

import (
	"log"
	"strings"

	"github.com/rudifa/go-github-access/pkg/ghaccess"
	"github.com/spf13/cobra"
)

// reposCmd represents the repo-list command
var reposCmd = &cobra.Command{
	Use:     "repo-list",
	Aliases: []string{"repos"},
	Short:   "Get repo list ",
	Long:    `Get repo list`,
	Run: func(cmd *cobra.Command, args []string) {
		user := cmd.Flag("user").Value.String()
		modeStr := cmd.Flag("mode").Value.String()
		mode, err := ghaccess.ParseMode(modeStr)
		if err != nil {
			log.Fatal(err)
		}
		ghaccess.GetRepos(user, mode)
	},
}

func init() {
	rootCmd.AddCommand(reposCmd)

	reposCmd.Flags().StringP("user", "u", "octocat", "github user name (default octocat)")
	reposCmd.Flags().StringP("mode", "m", "data", "output mode: "+strings.Join(ghaccess.ModeStrings(), ", "))
}
