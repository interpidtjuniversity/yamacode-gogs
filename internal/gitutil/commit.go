package gitutil

import (
	"bytes"
	"github.com/unknwon/com"
	"strings"
)

type bufferReader struct {
	buffer      []byte
	index       int
	currentLine []byte
	prevLine    []byte
	over        bool
}
func (br *bufferReader) readLine() {
	br.prevLine = make([]byte, len(br.currentLine))
	copy(br.prevLine, br.currentLine)
	br.currentLine = make([]byte, 0)
	for i:=br.index; br.index<len(br.buffer) && br.buffer[i]!='\n'; i++ {
		br.index++
		br.currentLine = append(br.currentLine, br.buffer[i])
	}
	if br.index >= len(br.buffer) {
		br.over = true
	} else {
		br.index++
	}
}

func GetRepoLatestMergeCommit(dir, branch string) (string, string, string, []byte, error) {
	stdout, stderr, err := com.ExecCmdDirBytes(dir, "git","log", branch)
	if err != nil {
		return "", "", "", stderr, err
	}
	commitInfo, mergeInfo := getLatestMergeCommitInfo(stdout)
	commitId := parseCommit(commitInfo)
	targetId, sourceId := parseMerge(mergeInfo)
	return commitId, sourceId, targetId, nil, nil
}

func getLatestMergeCommitInfo(buffer []byte) ([]byte, []byte){
	reader := bufferReader{buffer: buffer}

	for !reader.over && !bytes.HasPrefix(reader.currentLine, []byte("Merge")) {
		reader.readLine()
	}

	return reader.prevLine, reader.currentLine
}

func parseCommit(commit []byte) string{
	var commitId []byte
	if bytes.HasPrefix(commit, []byte("commit")) {
		commitId = commit[7:]
	}
	return string(commitId)
}

func parseMerge(buffer []byte) (string, string) {
	all := buffer[7:]
	allArray := strings.Split(string(all)," ")
	return allArray[0], allArray[1]
}
