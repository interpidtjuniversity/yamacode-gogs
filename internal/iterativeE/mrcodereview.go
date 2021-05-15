
package iterativeE

import (
	"github.com/gogs/git-module"
	"gogs.io/gogs/internal/conf"
	"gogs.io/gogs/internal/context"
	"gogs.io/gogs/internal/db"
	"gogs.io/gogs/internal/gitutil"
	invoker "gogs.io/gogs/internal/grpc/invoke/invokeImpl"
)

const (
	MR_FILES   = "repo/mrs/files"
)

type IterativeE struct{}

func GetDiff(sourceCommitId, targetCommitId string, gitRepo *git.Repository) (*gitutil.Diff, error){
	return gitutil.RepoDiff(gitRepo, sourceCommitId, conf.Git.MaxDiffFiles, conf.Git.MaxDiffLines, conf.Git.MaxDiffLineChars, git.DiffOptions{Base: targetCommitId})
}

func ViewMRFiles(c *context.Context) {
	c.Data["PageIsCodeReview"] = true
	index := c.ParamsInt64(":index")
	m, err := db.GetMergeRequestById(index)
	if err != nil {
		return
	}
	ownerName := c.Params("username")
	repoName := c.Params("reponame")
	_, err = db.GetUserByName(ownerName)
	if err != nil {
		return
	}
	gitRepo, err := git.Open(db.RepoPath(ownerName, repoName))
	if err != nil {
		return
	}

	diff, err := GetDiff(m.SourceCommitId, m.TargetCommitId, gitRepo)
	if err != nil {
		return
	}
	c.Data["MRDiff"] = diff
	c.Data["Reviewers"] = m.Reviewers
	c.Data["MRIndex"] = index

	c.Success(MR_FILES)
}

func PassMR(c *context.Context) {
	mrIndex := c.ParamsInt64(":index")
	user := c.Query("user")
	m, err := db.GetMergeRequestById(mrIndex)
	if err != nil {
		return
	}
	var remainViewers []string
	for _, v := range m.Reviewers {
		if v != user {
			remainViewers = append(remainViewers, v)
		}
	}
	db.UpdateMergeRequestViewersById(mrIndex, remainViewers)
	if len(remainViewers) == 0{
		// invoke grpc
		invoker.InvokePassMergerRequestCodeReview(m.ActionId, m.StageId, m.StepId)
	}
}
