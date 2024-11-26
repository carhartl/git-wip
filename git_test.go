package main

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestParseWithModified(t *testing.T) {
	var tests = []struct {
		name  string
		input string
	}{
		{"staged", "1 M. N... 100644 100644 100644 e69de29bb2d1d6434b8b29ae775ad8c2e48c5391 e69de29bb2d1d6434b8b29ae775ad8c2e48c5391 test.txt\n"},
		{"unstaged", "1 .M N... 100644 100644 100644 e69de29bb2d1d6434b8b29ae775ad8c2e48c5391 e69de29bb2d1d6434b8b29ae775ad8c2e48c5391 test.txt\n"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := strings.NewReader(tt.input)

			gi := GitInfo{}
			gi.Parse(s)

			require.Equal(t, 1, gi.modified)
		})
	}
}

func TestParseWithAdded(t *testing.T) {
	s := strings.NewReader("1 A. N... 000000 100644 100644 0000000000000000000000000000000000000000 e69de29bb2d1d6434b8b29ae775ad8c2e48c5391 test.txt\n")

	gi := GitInfo{}
	gi.Parse(s)

	require.Equal(t, 1, gi.added)
}

func TestParseWithDeleted(t *testing.T) {
	var tests = []struct {
		name  string
		input string
	}{
		{"staged", "1 D. N... 100644 000000 000000 e69de29bb2d1d6434b8b29ae775ad8c2e48c5391 0000000000000000000000000000000000000000 test.txt\n"},
		{"unstaged", "1 .D N... 100644 000000 000000 e69de29bb2d1d6434b8b29ae775ad8c2e48c5391 0000000000000000000000000000000000000000 test.txt\n"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := strings.NewReader(tt.input)

			gi := GitInfo{}
			gi.Parse(s)

			require.Equal(t, 1, gi.deleted)
		})
	}
}

func TestParseWithRenamed(t *testing.T) {
	s := strings.NewReader("2 R. N... 100644 100644 100644 e69de29bb2d1d6434b8b29ae775ad8c2e48c5391 e69de29bb2d1d6434b8b29ae775ad8c2e48c5391 R100 renamed.txt   test.txt\n")

	gi := GitInfo{}
	gi.Parse(s)

	require.Equal(t, 1, gi.renamed)
}

func TestParseWithCopied(t *testing.T) {
	s := strings.NewReader("2 CD N... 100644 100644 000000 e69de29bb2d1d6434b8b29ae775ad8c2e48c5391 e69de29bb2d1d6434b8b29ae775ad8c2e48c5391 C100 copied.txt test.txt")

	gi := GitInfo{}
	gi.Parse(s)

	require.Equal(t, 1, gi.copied)
}

func TestParseWithUnmerged(t *testing.T) {
	s := strings.NewReader("u UU N... 100644 100644 100644 100644 323fae03f4606ea9991df8befbb2fca795e648fa 257cc5642cb1a054f08cc83f2d943e56fd3ebe99 27f5cb292011032e79279cbd0fac3b1ecff8ce9a test.txt\n")

	gi := GitInfo{}
	gi.Parse(s)

	require.Equal(t, 1, gi.unmerged)
}

func TestParseWithUntracked(t *testing.T) {
	s := strings.NewReader("? test.txt\n")

	gi := GitInfo{}
	gi.Parse(s)

	require.Equal(t, 1, gi.untracked)
}

func TestParseWithStashed(t *testing.T) {
	s := strings.NewReader("# stash 1\n")

	gi := GitInfo{}
	gi.Parse(s)

	require.Equal(t, 1, gi.stashed)
}

func TestParseWithAhead(t *testing.T) {
	s := strings.NewReader("# branch.ab +1 -0\n")

	gi := GitInfo{}
	gi.Parse(s)

	require.Equal(t, 1, gi.ahead)
}

func TestIsClean(t *testing.T) {
	var gi GitInfo

	gi = GitInfo{}
	require.Equal(t, true, gi.IsClean())

	gi = GitInfo{modified: 1}
	require.Equal(t, false, gi.IsClean())
}

func TestSummary(t *testing.T) {
	var tests = []struct {
		name  string
		input GitInfo
		want  string
	}{
		{"withOneModified", GitInfo{modified: 1}, "1 files to commit"},
		{"withManyModified", GitInfo{modified: 1, added: 1, deleted: 1, renamed: 1, copied: 1, untracked: 1}, "6 files to commit"},
		{"withUnmerged", GitInfo{unmerged: 1}, "1 files to merge"},
		{"withStashed", GitInfo{stashed: 1}, "1 stashes"},
		{"withUnpushedCommits", GitInfo{ahead: 1}, "1 unpushed commits"},
		{"withManyDifferent", GitInfo{modified: 1, stashed: 1}, "1 files to commit, 1 stashes"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require.Equal(t, tt.want, tt.input.Summary())
		})
	}
}
