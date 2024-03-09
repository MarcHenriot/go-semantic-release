package git

import (
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/config"
	"github.com/go-git/go-git/v5/storage/memory"
	"github.com/hashicorp/go-version"
	"go.uber.org/zap"
)

type TagManager struct {
	logger *zap.Logger
	remote *git.Remote
}

type Tag struct {
	Name      string
	Version   *version.Version
	Ref       string
	CommitSha string
}

func NewTagManager(url string, logger *zap.Logger) *TagManager {
	logger.Info("creating tag manager", zap.String("url", url))
	remote := git.NewRemote(memory.NewStorage(), &config.RemoteConfig{
		Name: "origin",
		URLs: []string{url},
	})
	return &TagManager{
		logger: logger,
		remote: remote,
	}
}

func (tm *TagManager) GetTags() (tags []*Tag, err error) {
	tm.logger.Debug("getting tags")
	refs, err := tm.remote.List(&git.ListOptions{})
	if err != nil {
		tm.logger.Error("failed to get refs", zap.Error(err))
		return nil, err
	}
	for _, ref := range refs {
		if ref.Name().IsTag() {
			semver, err := version.NewSemver(ref.Name().Short())
			if err != nil {
				tm.logger.Warn("failed to parse semver", zap.Error(err))
			} else {
				tag := &Tag{
					Name:      ref.Name().Short(),
					Version:   semver,
					Ref:       ref.Name().String(),
					CommitSha: ref.String(),
				}
				tm.logger.Debug("found tag", zap.String("name", tag.Name), zap.String("ref", tag.Ref), zap.String("commitSha", tag.CommitSha))
				tags = append(tags, tag)
			}
		}
	}
	tm.logger.Info("GetTags", zap.Int("count", len(tags)))
	return tags, nil
}

func (tm *TagManager) GetLatestTag(tags []*Tag) (latestTag *Tag, err error) {
	tm.logger.Debug("getting latest tag")
	latestTag = tags[0]
	for _, tag := range tags {
		if tag.Version.GreaterThan(latestTag.Version) {
			latestTag = tag
		}
	}
	return latestTag, nil
}
