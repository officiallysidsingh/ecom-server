package services

import (
	"context"

	"github.com/officiallysidsingh/ecom-server/internal/models"
	"github.com/officiallysidsingh/ecom-server/internal/store"
)

type UserService interface {
	Login(ctx context.Context, loginReq *models.LoginRequest) (*models.User, error)
	Signup(ctx context.Context, user *models.User) (string, error)
}

type userService struct {
	store store.UserStore
}

func NewUserService(store store.UserStore) UserService {
	return &userService{
		store: store,
	}
}

func (s *userService) Login(ctx context.Context, loginReq *models.LoginRequest) (*models.User, error) {
	return s.store.GetByEmailFromDB(ctx, loginReq.Email)
}

func (s *userService) Signup(ctx context.Context, user *models.User) (string, error) {
	return s.store.CreateInDB(ctx, user)
}
