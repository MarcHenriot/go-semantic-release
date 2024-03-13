package git

import (
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/storage/memory"
	"go.uber.org/zap"
)

type RemoteManager struct {
	logger *zap.Logger
	remote *git.Remote
	repo   *git.Repository
}

func NewRemoteManager(url, branchName, remoteName string, logger *zap.Logger) (*RemoteManager, error) {
	logger.Info("creating tag manager", zap.String("url", url))
	repo, err := git.Clone(memory.NewStorage(), nil, &git.CloneOptions{
		URL:           url,
		RemoteName:    remoteName,
		ReferenceName: plumbing.NewBranchReferenceName(branchName),
		SingleBranch:  true,
		Tags:          git.AllTags,
	})
	if err != nil {
		logger.Error("failed to clone repo", zap.Error(err))
		return nil, err
	}
	remote, err := repo.Remote(remoteName)
	if err != nil {
		logger.Error("failed to get remote", zap.Error(err))
		return nil, err
	}
	return &RemoteManager{
		logger: logger,
		remote: remote,
		repo:   repo,
	}, nil
}

func (rm *RemoteManager) ListTags() (tags Tags, err error) {
	rm.logger.Debug("getting tags")
	refs, err := rm.remote.List(&git.ListOptions{})
	if err != nil {
		rm.logger.Error("failed to get refs", zap.Error(err))
		return nil, err
	}
	for _, ref := range refs {
		if ref.Name().IsTag() {
			tags = append(tags, NewTag(ref.Name().Short(), ref.Name().String(), ref.String()))
			rm.logger.Debug(
				"found tag",
				zap.String("name", ref.Name().Short()),
				zap.String("ref", ref.Name().String()),
				zap.String("commitSha", ref.String()),
			)
		}
	}
	rm.logger.Info("GetTags", zap.Int("count", len(tags)))
	return tags, nil
}

func (rm *RemoteManager) ListCommits() (tags Tags, err error) {
	refs, err := rm.remote.List(&git.ListOptions{})
	if err != nil {
		rm.logger.Error("failed to get refs", zap.Error(err))
		return nil, err
	}
	for _, ref := range refs {
		rm.logger.Info(ref.String())
	}
	return tags, nil
}
