package main

import (
	"fmt"

	"github.com/MarcHenriot/go-semantic-release/pkg/git"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

var logger *zap.Logger
var repoURL string

func init() {
	var err error
	logger, err = zap.NewProduction()
	if err != nil {
		fmt.Println(err)
	}
	logger.Info("logger initialized")

	// Add a persistent flag for the repository URL
	rootCmd.PersistentFlags().StringVarP(&repoURL, "repo", "r", "", "Repository URL (e.g., https://github.com/argoproj/argo-cd.git)")
	rootCmd.MarkPersistentFlagRequired("repo")
}

var rootCmd = &cobra.Command{
	Use:   "go-semantic-release",
	Short: "A CLI for semantic versioning and releasing based on commit and tag in CI",
	Long:  `This CLI automates semantic versioning and releases based on commit messages and tags in continuous integration (CI) environments.`,
}

var validateCmd = &cobra.Command{
	Use:   "validate",
	Short: "Validate commits for semantic versioning",
	Long:  `This command validates commits to ensure they adhere to semantic versioning conventions.`,
	Run: func(cmd *cobra.Command, args []string) {
	},
}

var releaseCmd = &cobra.Command{
	Use:   "release",
	Short: "Create a new release based on commit messages",
	Long:  `This command creates a new release based on commit messages and semantic versioning.`,
	Run: func(cmd *cobra.Command, args []string) {
		tm := git.NewTagManager(repoURL, logger)
		tags, err := tm.GetTags()
		if err != nil {
			logger.Error("Failed to get tags", zap.Error(err))
			return
		}
		latestTag, _ := tm.GetLatestTag(tags)
		logger.Info("Latest tag", zap.String("Name", latestTag.Name), zap.String("Ref", latestTag.Ref), zap.String("CommitSha", latestTag.CommitSha))
	},
}

func main() {
	rootCmd.AddCommand(validateCmd)
	rootCmd.AddCommand(releaseCmd)
	rootCmd.Execute()
}
