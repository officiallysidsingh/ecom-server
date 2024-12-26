package interfaces

import (
	"context"

	"github.com/officiallysidsingh/ecom-server/internal/models"
)

type ProductService interface {
	GetAll(ctx context.Context) ([]models.Product, error)
	GetByID(ctx context.Context, id string) (*models.Product, error)
	Create(ctx context.Context, product *models.Product) (string, error)
	PutUpdate(ctx context.Context, product *models.Product, id string) error
	PatchUpdate(ctx context.Context, product *models.Product, id string) error
	Delete(ctx context.Context, id string) error
}

type ProductStore interface {
	GetAllFromDB(ctx context.Context) ([]models.Product, error)
	GetByIDFromDB(ctx context.Context, id string) (*models.Product, error)
	CreateInDB(ctx context.Context, product *models.Product) (string, error)
	PutUpdateInDB(ctx context.Context, product *models.Product, id string) error
	PatchUpdateInDB(ctx context.Context, product *models.Product, id string) error
	DeleteFromDB(ctx context.Context, id string) error
}
