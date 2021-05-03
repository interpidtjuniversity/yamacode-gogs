package serviceImpl

import (
	"context"
	api "github.com/gogs/go-gogs-client"
	"gogs.io/gogs/internal/db"
	"gogs.io/gogs/internal/route/api/v1/repo"
)

type ApplicationService struct {
}

func (b *ApplicationService) QueryApplicationOwners(ctx context.Context, request *QueryOwnerRequest) (*QueryOwnerResponse, error) {
	response := &QueryOwnerResponse{}
	users, err := db.Users.GetAllUser()
	if err != nil {
		return response, err
	}
	response.OwnerNames = users
	return response, nil
}

func (b *ApplicationService) QueryApplications(ctx context.Context, request *QueryApplicationRequest) (*QueryApplicationResponse, error) {
	user, err := db.Users.GetByUsername(request.OwnerName)
	response := &QueryApplicationResponse{}
	if err != nil {
		return response, err
	}
	repos, err := db.Repos.GetByOwner(user.ID)
	if err != nil {
		return nil, err
	}
	response.AppNames = repos
	return response, nil
}

func (b *ApplicationService) CreateApplication(ctx context.Context, request *CreateApplicationRequest) (*CreateApplicationResponse, error) {
	// 1.query user
	baseResponse := &CreateApplicationResponse{Success: false}
	var user *db.User
	var err error
	if request.UserName != "" {
		user, err = db.Users.GetByUsername(request.UserName)
	} else {
		user, err = db.Users.GetByID(request.UserId)
	}
	if err != nil {
		return baseResponse, err
	}
	newRepo, err :=repo.CreateUserRepoWithoutContext(user, api.CreateRepoOption{
		Name:        request.RepoName,
		Description: request.Description,
		Private:     request.IsPrivate,
		AutoInit:    request.AutoInit,
	})
	if err != nil {
		return baseResponse, err
	}
	baseResponse.Success = true
	baseResponse.RepoName = newRepo.Name
	baseResponse.Owner = newRepo.Owner.UserName
	baseResponse.RepoId = newRepo.ID
	baseResponse.CloneUrl = newRepo.CloneURL
	baseResponse.SshUrl = newRepo.SSHURL
	baseResponse.HtmlUrl = newRepo.HTMLURL
	baseResponse.DefaultBranch = newRepo.DefaultBranch
	baseResponse.FullRepoName = newRepo.FullName
	baseResponse.Description = newRepo.Description
	baseResponse.WebSite = newRepo.Website
	baseResponse.Private = newRepo.Private
	return baseResponse, nil
}
