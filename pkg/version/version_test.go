package version

import (
	"testing"

	"github.com/MarcHenriot/go-semantic-release/pkg/git"
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
				git.NewTag("v1.3.0", "", ""),
				git.NewTag("2.2.3", "", ""),
				git.NewTag("1.3.1", "", ""),
				git.NewTag("v1.2.3", "", ""),
			},
			expectedTag: git.NewTag("2.2.3", "", ""),
		},
		{
			name: "test one bad",
			tags: git.Tags{
				git.NewTag("v1.3.0", "", ""),
				git.NewTag("sdfsdf", "", ""),
				git.NewTag("1.3.1", "", ""),
				git.NewTag("v1.2.3", "", ""),
			},
			expectedTag: git.NewTag("1.3.1", "", ""),
		},
		{
			name: "test two consecutive bad",
			tags: git.Tags{
				git.NewTag("sdfsdf", "", ""),
				git.NewTag("sdfsdf", "", ""),
				git.NewTag("1.3.1", "", ""),
				git.NewTag("v1.2.3", "", ""),
			},
			expectedTag: git.NewTag("1.3.1", "", ""),
		},
		{
			name: "test one good",
			tags: git.Tags{
				git.NewTag("sdfsdf", "", ""),
				git.NewTag("sdfsdf", "", ""),
				git.NewTag("1.3.1", "", ""),
				git.NewTag("sdfsdf", "", ""),
			},
			expectedTag: git.NewTag("1.3.1", "", ""),
		},
		{
			name: "test one bad at start",
			tags: git.Tags{
				git.NewTag("sdfsdf", "", ""),
				git.NewTag("0.2.3", "", ""),
				git.NewTag("1.3.1", "", ""),
				git.NewTag("8.9.45", "", ""),
			},
			expectedTag: git.NewTag("8.9.45", "", ""),
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
