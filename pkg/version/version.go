package version

import (
	"sort"

	"github.com/MarcHenriot/go-semantic-release/pkg/git"
	"github.com/Masterminds/semver/v3"
	"go.uber.org/zap"
)

type VersionManager struct {
	logger     *zap.Logger
	semverTags git.Tags
}

func NewVersionManager(tags git.Tags, logger *zap.Logger) *VersionManager {
	return &VersionManager{
		logger:     logger,
		semverTags: tags,
	}
}

func (vm *VersionManager) GetLatestTag() *git.Tag {
	if len(vm.semverTags) <= 0 {
		return nil
	}
	sort.Slice(vm.semverTags, func(i, j int) bool {
		semverI, errI := semver.NewVersion(vm.semverTags[i].GetName())
		if errI != nil {
			vm.logger.Warn(
				"failed to parse semver i",
				zap.String("Name", vm.semverTags[i].GetName()),
				zap.Error(errI),
			)
			return false
		}
		semverj, errJ := semver.NewVersion(vm.semverTags[j].GetName())
		if errJ != nil {
			vm.logger.Warn(
				"failed to parse semver j",
				zap.String("Name", vm.semverTags[j].GetName()),
				zap.Error(errJ),
			)
			return true
		}
		return semverI.GreaterThan(semverj)
	})
	return vm.semverTags[0]
}
