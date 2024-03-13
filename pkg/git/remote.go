package git

import (
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
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
			tags = append(tags, NewTag(ref.Name().Short(), ref.Name().String(), ref.Hash()))
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

func (rm *RemoteManager) ListCommitSince(tag *Tag) (commits Commits, err error) {
	commit, _ := rm.repo.CommitObject(tag.GetHash())
	tagTime := commit.Committer.When
	head, _ := rm.repo.Head()
	cIter, err := rm.repo.Log(&git.LogOptions{
		From:  head.Hash(),
		Since: &tagTime,
	})
	commits = make(Commits, 0)
	if err != nil {
		return commits, err
	}
	err = cIter.ForEach(func(c *object.Commit) error {
		rm.logger.Info("found commit", zap.String("message", c.Message))
		commits = append(commits, Commit(c.Message))
		return nil
	})
	return commits, err
}
