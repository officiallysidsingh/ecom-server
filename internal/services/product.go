package services

import (
	"context"
	"errors"
	"fmt"

	"github.com/officiallysidsingh/ecom-server/internal/interfaces"
	"github.com/officiallysidsingh/ecom-server/internal/models"
)

var (
	ErrInvalidPrice = errors.New("price must be greater than 0")
	ErrInvalidStock = errors.New("stock cannot be negative")
)

type ProductService struct {
	store interfaces.ProductStore
}

func NewProductService(store interfaces.ProductStore) interfaces.ProductService {
	return &ProductService{
		store: store,
	}
}

func (s *ProductService) GetAll(ctx context.Context) ([]models.Product, error) {
	return s.store.GetAllFromDB(ctx)
}

func (s *ProductService) GetByID(ctx context.Context, id string) (*models.Product, error) {
	return s.store.GetByIDFromDB(ctx, id)
}

func (s *ProductService) Create(ctx context.Context, product *models.Product) (string, error) {
	if product.Price <= 0 {
		return "", ErrInvalidPrice
	}
	if product.Stock < 0 {
		return "", ErrInvalidStock
	}

	return s.store.CreateInDB(ctx, product)
}

func (s *ProductService) PutUpdate(ctx context.Context, product *models.Product, id string) error {
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

func (s *ProductService) PatchUpdate(ctx context.Context, product *models.Product, id string) error {
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

func (s *ProductService) Delete(ctx context.Context, id string) error {
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
