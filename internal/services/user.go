package services

import (
	"context"

	"github.com/officiallysidsingh/ecom-server/internal/interfaces"
	"github.com/officiallysidsingh/ecom-server/internal/models"
)

type UserService struct {
	store interfaces.UserStore
}

func NewUserService(store interfaces.UserStore) interfaces.UserService {
	return &UserService{
		store: store,
	}
}

func (s *UserService) Login(ctx context.Context, loginReq *models.LoginRequest) (*models.User, error) {
	return s.store.GetByEmailFromDB(ctx, loginReq.Email)
}

func (s *UserService) Signup(ctx context.Context, user *models.User) (string, error) {
	return s.store.CreateInDB(ctx, user)
}
