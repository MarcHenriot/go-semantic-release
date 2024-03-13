package main

import (
	"fmt"

	"github.com/MarcHenriot/go-semantic-release/pkg/git"
	"github.com/MarcHenriot/go-semantic-release/pkg/version"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

var logger *zap.Logger
var repoURL string
var branchName string
var remoteName string

func init() {
	var err error
	logger, err = zap.NewProduction()
	if err != nil {
		fmt.Println(err)
	}
	logger.Debug("logger initialized")

	rootCmd.PersistentFlags().StringVarP(&repoURL, "repo-url", "u", "", "Repository URL (e.g., https://github.com/argoproj/argo-cd.git)")
	rootCmd.PersistentFlags().StringVarP(&branchName, "branch", "b", "main", "Branch name")
	rootCmd.PersistentFlags().StringVarP(&remoteName, "remote", "r", "origin", "Remote name")
	rootCmd.MarkPersistentFlagRequired("repo-url")
}

var rootCmd = &cobra.Command{
	Use:   "go-semantic-release",
	Short: "A CLI for semantic versioning and releasing based on commit and tag in CI",
	Long:  `This CLI automates semantic versioning and releases based on commit messages and tags in continuous integration (CI) environments.`,
	Args:  cobra.MinimumNArgs(1),
}

var validateCmd = &cobra.Command{
	Use:   "validate",
	Short: "Validate commits for semantic versioning",
	Long:  `This command validates commits to ensure they adhere to semantic versioning conventions.`,
	Run:   func(cmd *cobra.Command, args []string) {},
}

var releaseCmd = &cobra.Command{
	Use:   "release",
	Short: "Create a new release based on commit messages",
	Long:  `This command creates a new release based on commit messages and semantic versioning.`,
	Run: func(cmd *cobra.Command, args []string) {
		rm, err := git.NewRemoteManager(repoURL, branchName, remoteName, logger)
		if err != nil {
			logger.Error("Failed to create remote manager", zap.Error(err))
			return
		}
		tags, err := rm.ListTags()
		if err != nil {
			logger.Error("Failed to get tags", zap.Error(err))
			return
		}
		vm := version.NewVersionManager(tags, logger)
		latestTag := vm.GetLatestTag()
		logger.Info(
			"Latest tag",
			zap.String("Name", latestTag.GetName()),
			zap.String("Ref", latestTag.GetRef()),
		)
		commits, _ := rm.ListCommitSince(latestTag)
		logger.Info("Commits since last tag", zap.Int("Count", len(commits)))
		// semverManager := semver.NewSemverManager(latestTag)
		// semverManager.NextMajor()
		// fmt.Println(semverManager.GetTag())
	},
}

func main() {
	rootCmd.AddCommand(validateCmd)
	rootCmd.AddCommand(releaseCmd)
	rootCmd.Execute()
}
