package services

import (
	"context"
	"errors"
	"fmt"

	"github.com/officiallysidsingh/ecom-server/internal/models"
	"github.com/officiallysidsingh/ecom-server/internal/store"
)

var (
	ErrInvalidPrice = errors.New("price must be greater than 0")
	ErrInvalidStock = errors.New("stock cannot be negative")
)

type ProductService interface {
	GetAll(ctx context.Context) ([]models.Product, error)
	GetByID(ctx context.Context, id string) (*models.Product, error)
	Create(ctx context.Context, product *models.Product) (string, error)
	PutUpdate(ctx context.Context, product *models.Product, id string) error
	PatchUpdate(ctx context.Context, product *models.Product, id string) error
	Delete(ctx context.Context, id string) error
}

type productService struct {
	store store.ProductStore
}

func NewProductService(store store.ProductStore) ProductService {
	return &productService{
		store: store,
	}
}

func (s *productService) GetAll(ctx context.Context) ([]models.Product, error) {
	return s.store.GetAllFromDB(ctx)
}

func (s *productService) GetByID(ctx context.Context, id string) (*models.Product, error) {
	return s.store.GetByIDFromDB(ctx, id)
}

func (s *productService) Create(ctx context.Context, product *models.Product) (string, error) {
	if product.Price <= 0 {
		return "", ErrInvalidPrice
	}
	if product.Stock < 0 {
		return "", ErrInvalidStock
	}

	return s.store.CreateInDB(ctx, product)
}

func (s *productService) PutUpdate(ctx context.Context, product *models.Product, id string) error {
	if product.Price <= 0 {
		return ErrInvalidPrice
	}
	if product.Stock < 0 {
		return ErrInvalidStock
	}

	_, err := s.store.GetByIDFromDB(ctx, id)
	if err != nil {
		return err
	}

	err = s.store.PutUpdateInDB(ctx, product, id)
	if err != nil {
		return err
	}

	return nil
}

func (s *productService) PatchUpdate(ctx context.Context, product *models.Product, id string) error {
	if product.Price != 0 && product.Price <= 0 {
		return ErrInvalidPrice
	}
	if product.Stock < 0 {
		return ErrInvalidStock
	}

	_, err := s.store.GetByIDFromDB(ctx, id)
	if err != nil {
		return err
	}

	err = s.store.PatchUpdateInDB(ctx, product, id)
	if err != nil {
		return err
	}

	return nil
}

func (s *productService) Delete(ctx context.Context, id string) error {
	_, err := s.store.GetByIDFromDB(ctx, id)
	if err != nil {
		return err
	}

	err = s.store.DeleteFromDB(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to delete product with ID %s: %v", id, err)
	}

	return nil
}
