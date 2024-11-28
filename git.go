package main

import (
	"bufio"
	"fmt"
	"io"
	"strconv"
	"strings"
)

type Dir = string

type GitInfo struct {
	modified    int
	typeChanged int
	added       int
	deleted     int
	renamed     int
	copied      int
	unmerged    int
	untracked   int
	stashed     int
	ahead       int
}

func (gi *GitInfo) Parse(r io.Reader) {
	var s = bufio.NewScanner(r)
	for s.Scan() {
		gi.parseLine(s.Text())
	}
}

func (gi GitInfo) IsClean() bool {
	return gi.modified == 0 &&
		gi.typeChanged == 0 &&
		gi.added == 0 &&
		gi.deleted == 0 &&
		gi.renamed == 0 &&
		gi.copied == 0 &&
		gi.unmerged == 0 &&
		gi.untracked == 0 &&
		gi.stashed == 0 &&
		gi.ahead == 0
}

func (gi GitInfo) Summary() string {
	s := []string{}

	committable := gi.modified +
		gi.typeChanged +
		gi.added +
		gi.deleted +
		gi.renamed +
		gi.copied +
		gi.untracked
	if committable > 0 {
		s = append(s, fmt.Sprintf("%d files to commit", committable))
	}

	if gi.unmerged > 0 {
		s = append(s, fmt.Sprintf("%d files to merge", gi.unmerged))
	}

	if gi.stashed > 0 {
		s = append(s, fmt.Sprintf("%d stashes", gi.stashed))
	}

	if gi.ahead > 0 {
		s = append(s, fmt.Sprintf("%d unpushed commits", gi.ahead))
	}

	return strings.Join(s, ", ")
}

func (gi *GitInfo) parseLine(l string) {
	s := bufio.NewScanner(strings.NewReader(l))
	s.Split(bufio.ScanWords)
	s.Scan()
	switch s.Text() {
	case "#":
		gi.parseHeader(l)
	case "1", "2":
		s.Scan()
		gi.parseXY(s.Text())
	case "u":
		gi.unmerged++
	case "?":
		gi.untracked++
	}
}

func (gi *GitInfo) parseXY(xy string) {
	// x: staged, y: unstaged
	for _, c := range xy {
		switch c { // staged
		case 'M':
			gi.modified++
		case 'T':
			gi.typeChanged++
		case 'A':
			gi.added++
		case 'D':
			gi.deleted++
		case 'R':
			gi.renamed++
		case 'C':
			gi.copied++
		}
	}
}

func (gi *GitInfo) parseHeader(s string) {
	parts := strings.Split(s, " ")
	switch parts[1] {
	case "stash": // line: # stash 1
		n, _ := strconv.Atoi(parts[2])
		gi.stashed = n
	case "branch.ab": // line: # branch.ab +1 -0
		n, _ := strconv.Atoi(parts[2])
		gi.ahead = n
	}
}
