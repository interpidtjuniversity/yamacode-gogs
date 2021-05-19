package db

import "xorm.io/builder"

type MergeRequestPipeline struct {
	ID              int64    `xorm:"id autoincr pk"`
	UserName        string   `xorm:"user_name"`
	RepoName        string   `xorm:"repo_name"`
	ActionId        int64    `xorm:"action_id"`
	SourceBranch    string   `xorm:"source_branch"`
	TargetBranch    string   `xorm:"target_branch"`
	Finish          bool     `xorm:"finish"`
	PipelineId      int64    `xorm:"pipeline_id"`
	PusherName      string   `xorm:"pusher_name"`
	IterationId     int64    `xorm:"iteration_id"`
	Env             string   `xorm:"env"`
	ActionInfo      string   `xorm:"action_info"`
	MRCodeReviewers []string `xorm:"reviewers"`
	MRInfo          string   `xorm:"mr_info"`
}

func InsertMergeRequestPipeline(mrp *MergeRequestPipeline) (int64, error){
	_,err := x.Insert(mrp)
	if err != nil {
		return 0, err
	}
	return mrp.ID, nil
}

func GetMergeRequestPipeline(userName, repoName, branch, pusherName string) ([]*MergeRequestPipeline, error){
	var mrps []*MergeRequestPipeline
	err := x.Table("merge_request_pipeline").Where(builder.Eq{"user_name":userName,
		"repo_name": repoName, "pusher_name":pusherName, "source_branch":branch,
		}).Find(&mrps)

	return mrps, err
}

func BranchUpdateMergeRequestPipelineFinish(actionIds []int64) error {
	mrp := &MergeRequestPipeline{Finish: true}
	_, err := x.Table("merge_request_pipeline").Cols("finish").Where(builder.In("action_id", actionIds)).Update(mrp)
	return err
}

