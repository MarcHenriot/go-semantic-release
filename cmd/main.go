package main

import (
	"fmt"
	"sort"

	"github.com/go-git/go-billy/v5/memfs"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/storage/memory"
	"github.com/hashicorp/go-version"
	"go.uber.org/zap"
)

func main() {
	logger, _ := zap.NewProduction()
	defer logger.Sync()
	fs := memfs.New()
	storer := memory.NewStorage()
	r, err := git.Clone(storer, fs, &git.CloneOptions{
		URL: "https://github.com/argoproj/argo-cd",
	})
	if err != nil {
		logger.Error("Error cloning", zap.Error(err))
	}
	tags, err := r.Tags()
	if err != nil {
		logger.Error("Error fetching tags", zap.Error(err))
	}
	var tagList []*version.Version
	err = tags.ForEach(func(t *plumbing.Reference) error {
		tag_version, err := version.NewSemver(t.Name().Short())
		if err != nil {
			logger.Warn("Cannot parse tag", zap.Error(err))
		} else {
			tagList = append(tagList, tag_version)
		}
		return nil
	})
	if err != nil {
		logger.Error("Error parsing tags", zap.Error(err))
	}
	sort.Sort(version.Collection(tagList))
	for _, tag := range tagList {
		fmt.Println(tag.String())
	}
}
