package db

import (
	"fmt"
	"gogs.io/gogs/internal/errutil"
)

type MergeRequest struct {
	ID             int64    `xorm:"id autoincr pk"`
	RepoId         int64    `xorm:"repo_id"`
	RepoName       string   `xorm:"repo_name"`
	SourceBranch   string   `xorm:"source_branch"`
	TargetBranch   string   `xorm:"target_branch"`
	SourceCommitId string   `xorm:"source_commit"`
	TargetCommitId string   `xorm:"target_commit"`
}

func InsertMergeRequest(m *MergeRequest) (int64, error) {
	_,err := x.Insert(m)
	if err != nil {
		return 0, err
	}
	return m.ID, nil
}

func GetMergeRequestById(index int64) (*MergeRequest, error) {
	mr := &MergeRequest{}
	has, _ := x.ID(index).Get(mr)
	if !has{
		return nil, ErrMRNotExist{Args: errutil.Args{"MRId": index}}
	}
	return mr, nil
}


type ErrMRNotExist struct {
	Args map[string]interface{}
}

func (err ErrMRNotExist) Error() string {
	return fmt.Sprintf("mergeRequest does not exist: %v", err.Args)
}

func (ErrMRNotExist) NotFound() bool {
	return true
}
