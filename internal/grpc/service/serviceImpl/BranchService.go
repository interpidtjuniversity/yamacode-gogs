package serviceImpl

import (
	"context"
	"fmt"
	"github.com/gogs/git-module"
	"gogs.io/gogs/internal/conf"
	"gogs.io/gogs/internal/db"
	"gogs.io/gogs/internal/errutil"
	"gogs.io/gogs/internal/gitutil"
)

type BranchService struct {}

func (b *BranchService) CreateBranch(ctx context.Context, request *CreateBranchRequest) (*CreateBranchResponse, error) {
	return nil,nil
}

func (b *BranchService) DeleteBranch(ctx context.Context, request *DeleteBranchRequest) (*DeleteBranchResponse, error) {
	return nil,nil
}

func (b *BranchService) Merge2Branch(ctx context.Context, request *Merge2BranchRequest) (*Merge2BranchResponse, error) {
	return nil,nil
}

func (b *BranchService) Query2BranchConflict(ctx context.Context, request *ConflictDetectRequest) (*ConflictDetectResponse, error) {
	return nil, nil
}

func (b *BranchService) QueryRepoBranchCommit(ctx context.Context, request *CommitQueryRequest) (*CommitQueryResponse, error) {
	baseResponse := &CommitQueryResponse{}
	ownerName := request.OwnerName
	repoName := request.RepoName
	owner, err := db.GetUserByName(ownerName)
	if err != nil {
		return baseResponse, db.ErrUserNotExist{Args: errutil.Args{"userName": ownerName}}
	}
	repository, err := db.GetRepositoryByName(owner.ID, repoName)
	if err != nil {
		return baseResponse, db.ErrRepoNotExist{Args: errutil.Args{"userName": ownerName, "repoName": repoName}}
	}
	gitRepo, err := git.Open(db.RepoPath(ownerName, repoName))
	if err != nil {
		return baseResponse, db.ErrRepoNotExist{Args: errutil.Args{"ownerName": ownerName, "repoName": repoName}}
	}
	/**
	               latestCommitId
	          parent             parent
	     sourcecCommitId    targetCommitId
	*/
	latestCommitId,sourcecCommitId,targetCommitId,_,_ := gitutil.GetRepoLatestMergeCommit(gitRepo.Path())
	mr, _ := db.GetMergeRequestByRepoIdAndCommitId(repository.ID, sourcecCommitId, targetCommitId)

	baseResponse.CommitId = latestCommitId

	var listenAddr string
	if conf.Server.Protocol == "unix" {
		listenAddr = conf.Server.HTTPAddr
	} else {
		listenAddr = fmt.Sprintf("%s:%s", conf.Server.HTTPAddr, conf.Server.HTTPPort)
	}
	baseResponse.Url = fmt.Sprintf("%v://%s%s/%s/%s/mr/%d/files", conf.Server.Protocol, listenAddr, conf.Server.Subpath, ownerName, repoName, mr.ID)
	//baseResponse.Url = fmt.Sprintf("%v://%s%s/%s/%s/mr/%d/files", conf.Server.Protocol, "localhost:3002", conf.Server.Subpath, ownerName, repoName, mr.ID)

	return baseResponse, nil
}

func (b *BranchService) RegisterMergeRequest(ctx context.Context, request *RegisterMRRequest) (*RegisterMRResponse, error) {
	baseResponse := &RegisterMRResponse{}

	ownerName := request.OwnerName
	repoName := request.RepoName
	owner, err := db.GetUserByName(ownerName)
	if err != nil {
		return baseResponse, db.ErrUserNotExist{Args: errutil.Args{"userName": ownerName}}
	}
	repository, err := db.GetRepositoryByName(owner.ID, repoName)
	if err != nil {
		return baseResponse, db.ErrRepoNotExist{Args: errutil.Args{"userName": ownerName, "repoName": repoName}}
	}
	gitRepo, err := git.Open(db.RepoPath(ownerName, repoName))
	if err != nil {
		return baseResponse, db.ErrRepoNotExist{Args: errutil.Args{"ownerName": ownerName, "repoName": repoName}}
	}
	sourceCommit, _ := gitRepo.BranchCommit(request.SourceBranch)
	targetCommit, _ := gitRepo.BranchCommit(request.TargetBranch)

	mr := db.MergeRequest{
		RepoId: repository.ID,
		RepoName: repoName,
		SourceBranch: request.SourceBranch,
		TargetBranch: request.TargetBranch,
		SourceCommitId: sourceCommit.ID.String(),
		TargetCommitId: targetCommit.ID.String(),
	}
	// insert this mr
	id, err := db.InsertMergeRequest(&mr)
	if err != nil {
		return baseResponse, err
	}

	var listenAddr string
	if conf.Server.Protocol == "unix" {
		listenAddr = conf.Server.HTTPAddr
	} else {
		listenAddr = fmt.Sprintf("%s:%s", conf.Server.HTTPAddr, conf.Server.HTTPPort)
	}
	baseResponse.MRId = id
	baseResponse.ShowDiffUri = fmt.Sprintf("%v://%s%s/%s/%s/mr/%d/files", conf.Server.Protocol, listenAddr, conf.Server.Subpath, ownerName, repoName, id)
	//baseResponse.ShowDiffUri = fmt.Sprintf("%v://%s%s/%s/%s/mr/%d/files", conf.Server.Protocol, "localhost:3002", conf.Server.Subpath, ownerName, repoName, id)
	return baseResponse, nil
}

