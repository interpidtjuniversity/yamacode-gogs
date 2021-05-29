package db

import (
	"fmt"
	"gogs.io/gogs/internal/errutil"
	"xorm.io/builder"
)

type MergeRequest struct {
	ID             int64    `xorm:"id autoincr pk"`
	RepoId         int64    `xorm:"repo_id"`
	RepoName       string   `xorm:"repo_name"`
	SourceBranch   string   `xorm:"source_branch"`
	TargetBranch   string   `xorm:"target_branch"`
	SourceCommitId string   `xorm:"source_commit"`
	TargetCommitId string   `xorm:"target_commit"`
	ActionId       int64    `xorm:"action_id"`
	StageId        int64    `xorm:"stage_id"`
	StepId         int64    `xorm:"step_id"`
	Reviewers      []string `xorm:"reviewers"`
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

func GetMergeRequestByRepoIdAndCommitId(repoId int64, sourceCommit, targetCommit string) (*MergeRequest, error) {
	mr := &MergeRequest{}
	has, err := x.Where(builder.Eq{"repo_id":repoId}.And(builder.Like{"source_commit", sourceCommit}).And(builder.Like{"target_commit", targetCommit})).Get(mr)
	if !has {
		return nil, err
	}
	return mr, nil
}

func GetLatestMergeRequestByTargetBranch(repoId int64, branchName string) (*MergeRequest, error) {
	mr := &MergeRequest{}
	_, err := x.Where(builder.Eq{"repo_id": repoId, "target_branch": branchName}).Desc("id").Limit(1).Get(mr)
	return mr, err
}

func UpdateMergeRequestViewersById(index int64, remainViewers []string) error {
	mr := &MergeRequest{Reviewers: remainViewers}
	_, err := x.Table("merge_request").Cols("reviewers").Where(builder.Eq{"id":index}).Update(mr)
	return err
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
