package gitutil

import (
	"fmt"
	"github.com/unknwon/com"
	"io/ioutil"
)

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
