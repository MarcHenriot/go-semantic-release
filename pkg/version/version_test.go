package version

import (
	"testing"

	"github.com/MarcHenriot/go-semantic-release/pkg/git"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func Test_GetLatestTag(t *testing.T) {
	tests := []struct {
		name        string
		tags        git.Tags
		expectedTag *git.Tag
	}{
		{
			name:        "test empty",
			tags:        git.Tags{},
			expectedTag: nil,
		},
		{
			name: "test all good with v",
			tags: git.Tags{
				git.NewTag("v1.3.0", "", plumbing.Hash{}),
				git.NewTag("2.2.3", "", plumbing.Hash{}),
				git.NewTag("1.3.1", "", plumbing.Hash{}),
				git.NewTag("v1.2.3", "", plumbing.Hash{}),
			},
			expectedTag: git.NewTag("2.2.3", "", plumbing.Hash{}),
		},
		{
			name: "test one bad",
			tags: git.Tags{
				git.NewTag("v1.3.0", "", plumbing.Hash{}),
				git.NewTag("sdfsdf", "", plumbing.Hash{}),
				git.NewTag("1.3.1", "", plumbing.Hash{}),
				git.NewTag("v1.2.3", "", plumbing.Hash{}),
			},
			expectedTag: git.NewTag("1.3.1", "", plumbing.Hash{}),
		},
		{
			name: "test two consecutive bad",
			tags: git.Tags{
				git.NewTag("sdfsdf", "", plumbing.Hash{}),
				git.NewTag("sdfsdf", "", plumbing.Hash{}),
				git.NewTag("1.3.1", "", plumbing.Hash{}),
				git.NewTag("v1.2.3", "", plumbing.Hash{}),
			},
			expectedTag: git.NewTag("1.3.1", "", plumbing.Hash{}),
		},
		{
			name: "test one good",
			tags: git.Tags{
				git.NewTag("sdfsdf", "", plumbing.Hash{}),
				git.NewTag("sdfsdf", "", plumbing.Hash{}),
				git.NewTag("1.3.1", "", plumbing.Hash{}),
				git.NewTag("sdfsdf", "", plumbing.Hash{}),
			},
			expectedTag: git.NewTag("1.3.1", "", plumbing.Hash{}),
		},
		{
			name: "test one bad at start",
			tags: git.Tags{
				git.NewTag("sdfsdf", "", plumbing.Hash{}),
				git.NewTag("0.2.3", "", plumbing.Hash{}),
				git.NewTag("1.3.1", "", plumbing.Hash{}),
				git.NewTag("8.9.45", "", plumbing.Hash{}),
			},
			expectedTag: git.NewTag("8.9.45", "", plumbing.Hash{}),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			vm := NewVersionManager(test.tags, zap.NewExample())
			actualTag := vm.GetLatestTag()
			assert.Equal(t, test.expectedTag, actualTag)
		})
	}
}
