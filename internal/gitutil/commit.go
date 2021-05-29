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
	firstLine   []byte
	over        bool
	findFirst   bool
}
func (br *bufferReader) readLine() {
	br.prevLine = make([]byte, len(br.currentLine))
	copy(br.prevLine, br.currentLine)
	if !br.findFirst {
		br.firstLine = make([]byte, len(br.currentLine))
		copy(br.firstLine, br.currentLine)
		if len(br.firstLine) > 0 {
			br.findFirst = true
		}
	}
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

func GetRepoLatestMergeCommit(dir, branch string) (string, []byte, error) {
	stdout, stderr, err := com.ExecCmdDirBytes(dir, "git","log", branch)
	if err != nil {
		return "", stderr, err
	}
	LatestCommitInfo, _, _ := getLatestMergeCommitInfo(stdout)
	LatestCommitId := parseCommit(LatestCommitInfo)
	//targetId, sourceId := parseMerge(mergeInfo)
	// sourceId, targetId := parseMerge(mergeInfo)
	return LatestCommitId , nil, nil
}

func getLatestMergeCommitInfo(buffer []byte) ([]byte, []byte, []byte){
	reader := bufferReader{buffer: buffer}

	for !reader.over && !bytes.HasPrefix(reader.currentLine, []byte("Merge")) {
		reader.readLine()
	}

	return reader.firstLine, reader.prevLine, reader.currentLine
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
