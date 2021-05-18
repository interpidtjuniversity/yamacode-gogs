package gitutil

import (
	"context"
	"fmt"
	"github.com/unknwon/com"
	"gogs.io/gogs/internal/strutil"
	"io/ioutil"
	"os"
	"os/exec"
)

const GITPUSH_DIR = "/var/run/yama/yamaHub/repo/%s"
const GITPUSH_LOG_DIR = "/var/run/yama/yamaHub/repo/%s/log.log"

// cd /refs/heads and copy master as newBranch
func NewRepoBranch(dir, newBranch string) (bool, error){
	headsDir := fmt.Sprintf("%s/refs/heads", dir)
	heads, _ := ioutil.ReadDir(headsDir)
	for _, head := range heads{
		if newBranch == head.Name() {
			return false, fmt.Errorf("error while create new branch, %s already exist", newBranch)
		}
	}
	_,bufErr,err := com.ExecCmdDirBytes(headsDir, "cp","master", newBranch)
	if err!=nil || len(bufErr) > 0 {
		return false, err
	}
	return true, nil

}

func ListRepoBranch(dir string) []string {
	var branches []string
	headsDir := fmt.Sprintf("%s/refs/heads", dir)
	heads, _ := ioutil.ReadDir(headsDir)

	for _, head := range heads {
		branches = append(branches, head.Name())
	}
	return branches
}

// git clone repository
// git checkout target
// git merge source
// git push repoPath target
func MergeSourceToTarget(dir, repoName, source, target, mergeInfo string) (bool, error) {
	ctx := context.Background()
	rdmDir, _ := strutil.RandomChars(10)
	yamaHubDir, _ := os.Getwd()
	mergeSH := fmt.Sprintf("%s/internal/gitutil/merge.sh", yamaHubDir)

	mergeCmd := exec.CommandContext(ctx, mergeSH, dir, target, source, repoName, mergeInfo)
	mergeCmd.Dir = fmt.Sprintf(GITPUSH_DIR, rdmDir)
	os.MkdirAll(mergeCmd.Dir, os.ModePerm)
	log, _ := os.OpenFile(fmt.Sprintf(GITPUSH_LOG_DIR, rdmDir), os.O_CREATE|os.O_WRONLY, 0777)
	mergeCmd.Stdout = log
	mergeCmd.Stderr = log
	err := mergeCmd.Run()
	if err != nil {
		return false, err
	}
	return true, err
}
