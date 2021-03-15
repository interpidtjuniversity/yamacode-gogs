package serviceImpl

import (
	"context"
	"fmt"
	"github.com/gogs/git-module"
	"gogs.io/gogs/internal/conf"
	"gogs.io/gogs/internal/db"
	"gogs.io/gogs/internal/errutil"
	"gogs.io/gogs/internal/grpc/service"
	"google.golang.org/grpc"
	"log"
	"net"
)

type BranchService struct {}

func (b *BranchService) CreateBranch(ctx context.Context, request *service.CreateBranchRequest) (*service.CreateBranchResponse, error) {
	return nil,nil
}

func (b *BranchService) DeleteBranch(ctx context.Context, request *service.DeleteBranchRequest) (*service.DeleteBranchResponse, error) {
	return nil,nil
}

func (b *BranchService) Merge2Branch(ctx context.Context, request *service.Merge2BranchRequest) (*service.Merge2BranchResponse, error) {
	return nil,nil
}

func (b *BranchService) Query2BranchConflict(ctx context.Context, request *service.ConflictDetectRequest) (*service.ConflictDetectResponse, error) {
	return nil, nil
}

func (b *BranchService) RegisterMergeRequest(ctx context.Context, request *service.RegisterMRRequest) (*service.RegisterMRResponse, error) {
	baseResponse := &service.RegisterMRResponse{}

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
	return baseResponse, nil
}

func StartBranchService() {
	lis, err := net.Listen("tcp", ":8000")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	//建立 gPRC 服务器，并注册服务
	s := grpc.NewServer()
	service.RegisterYaMaHubBranchServiceServer(s, &BranchService{})

	log.Println("Server run ...")
	//启动服务
	if err := s.Serve(lis); err != nil {
		log.Fatalf("fail to serve: %v", err)
	}
}

