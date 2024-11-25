package main

import (
	"bufio"
	"io"
	"strconv"
	"strings"
)

type Dir = string

type GitInfo struct {
	modified  int
	added     int
	deleted   int
	renamed   int
	copied    int
	unmerged  int
	untracked int
	stashed   int
}

func (gi *GitInfo) Parse(r io.Reader) {
	var s = bufio.NewScanner(r)
	for s.Scan() {
		gi.parseLine(s.Text())
	}
}

func (gi GitInfo) Summary() string {
	// TODO: produces human readable output
	return ""
}

func (gi GitInfo) IsClean() bool {
	return gi.modified == 0 &&
		gi.added == 0 &&
		gi.deleted == 0 &&
		gi.renamed == 0 &&
		gi.copied == 0 &&
		gi.unmerged == 0 &&
		gi.untracked == 0 &&
		gi.stashed == 0
}

func (gi *GitInfo) parseLine(l string) {
	s := bufio.NewScanner(strings.NewReader(l))
	s.Split(bufio.ScanWords)
	s.Scan()
	switch s.Text() {
	case "#":
		gi.parseStashes(l)
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

func (gi *GitInfo) parseStashes(s string) {
	// line: # stash <N>
	stashed := strings.Split(s, " ")
	if stashed[1] == "stash" {
		n, _ := strconv.Atoi(stashed[2])
		gi.stashed = n
	}
}
