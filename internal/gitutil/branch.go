package gitutil

import "github.com/unknwon/com"

func NewRepoBranch(dir, newBranch string) (bool, error){
	// checkout master
	_,_,err := com.ExecCmdDirBytes(dir, "checkout","master")
	if err != nil {
		return false, err
	}
	// checkout -b newbranch
	_,_,err = com.ExecCmdDirBytes(dir, "checkout","-b", newBranch)
	if err != nil {
		return false, err
	}

	return true, nil
}
