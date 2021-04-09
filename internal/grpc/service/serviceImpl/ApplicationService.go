package serviceImpl

import (
	"context"
	"gogs.io/gogs/internal/db"
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
