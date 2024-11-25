package main

import (
	"strings"
	"testing"
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

			if gi.modified != 1 {
				t.Errorf("Expected modified == 1, got: %v", gi.modified)
			}
		})
	}
}

func TestParseWithAdded(t *testing.T) {
	s := strings.NewReader("1 A. N... 000000 100644 100644 0000000000000000000000000000000000000000 e69de29bb2d1d6434b8b29ae775ad8c2e48c5391 test.txt\n")

	gi := GitInfo{}
	gi.Parse(s)

	if gi.added != 1 {
		t.Errorf("Expected added == 1, got: %v", gi.added)
	}
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

			if gi.deleted != 1 {
				t.Errorf("Expected deleted == 1, got: %v", gi.deleted)
			}
		})
	}
}

func TestParseWithRenamed(t *testing.T) {
	s := strings.NewReader("2 R. N... 100644 100644 100644 e69de29bb2d1d6434b8b29ae775ad8c2e48c5391 e69de29bb2d1d6434b8b29ae775ad8c2e48c5391 R100 renamed.txt   test.txt\n")

	gi := GitInfo{}
	gi.Parse(s)

	if gi.renamed != 1 {
		t.Errorf("Expected renamed == 1, got: %v", gi.renamed)
	}
}

func TestParseWithCopied(t *testing.T) {
	s := strings.NewReader("2 CD N... 100644 100644 000000 e69de29bb2d1d6434b8b29ae775ad8c2e48c5391 e69de29bb2d1d6434b8b29ae775ad8c2e48c5391 C100 copied.txt test.txt")

	gi := GitInfo{}
	gi.Parse(s)

	if gi.copied != 1 {
		t.Errorf("Expected copied == 1, got: %v", gi.copied)
	}
}

func TestParseWithUnmerged(t *testing.T) {
	s := strings.NewReader("u UU N... 100644 100644 100644 100644 323fae03f4606ea9991df8befbb2fca795e648fa 257cc5642cb1a054f08cc83f2d943e56fd3ebe99 27f5cb292011032e79279cbd0fac3b1ecff8ce9a test.txt\n")

	gi := GitInfo{}
	gi.Parse(s)

	if gi.unmerged != 1 {
		t.Errorf("Expected unmerged == 1, got: %v", gi.unmerged)
	}
}

func TestParseWithUntracked(t *testing.T) {
	s := strings.NewReader("? test.txt\n")

	gi := GitInfo{}
	gi.Parse(s)

	if gi.untracked != 1 {
		t.Errorf("Expected untracked == 1, got: %v", gi.untracked)
	}
}

func TestParseWithStashed(t *testing.T) {
	s := strings.NewReader("# stash 1\n")

	gi := GitInfo{}
	gi.Parse(s)

	if gi.stashed != 1 {
		t.Errorf("Expected stashed == 1, got: %v", gi.stashed)
	}
}
