package interfaces

import (
	"context"

	"github.com/officiallysidsingh/ecom-server/internal/models"
)

type UserService interface {
	Login(ctx context.Context, email string, password string) (*models.User, error)
	Signup(ctx context.Context, user *models.User) (string, error)
}

type UserStore interface {
	GetByEmailFromDB(ctx context.Context) ([]models.User, error)
	CreateInDB(ctx context.Context, user *models.User) (string, error)
}
