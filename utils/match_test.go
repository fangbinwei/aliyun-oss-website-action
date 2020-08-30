package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type mock struct {
	pattern string
	ossPath string
	expect  bool
}

func TestMatch(t *testing.T) {
	assert := assert.New(t)

	matches := []mock{
		{"dir/", "dir", false},
		{"dir/", "dir/sub", true},
		{"./dir/", "dir/file", true},
		{"dir/", "dir/dir2/file", true},
		{"dir2/", "dir/dir2/file", false},

		{"file", "file", true},
		{"file", "file1", false},
		{"file*", "file1", true},
		{"file", "file/", false},
		{"file", "dir/file", false},
		{"dir/*.js", "dir/file.js", true},
		{"dir*/", "dir1/", true},
		{"dir/*.js", "dir/file.css", false},
		{"dir/*.js", "dir/dir2/file.js", false},
		{"dir/*/*.js", "dir/dir2/file.js", true},
		{"dir/**/*.js", "dir/dir2/dir3/file.js", false},
	}

	for _, item := range matches {
		if item.expect {
			assert.True(match(item.pattern, item.ossPath), item)
		} else {
			assert.False(match(item.pattern, item.ossPath), item)
		}
	}
}
