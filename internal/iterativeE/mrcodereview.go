
package iterativeE

import (
	"github.com/gogs/git-module"
	"gogs.io/gogs/internal/conf"
	"gogs.io/gogs/internal/context"
	"gogs.io/gogs/internal/db"
	"gogs.io/gogs/internal/gitutil"
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
	m, err := db.GetMergeRequestById(c.ParamsInt64(":index"))
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

	c.Success(MR_FILES)
}
