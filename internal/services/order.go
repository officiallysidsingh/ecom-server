package services

import (
	"context"
	"errors"
	"fmt"
	"log"

	"github.com/officiallysidsingh/ecom-server/internal/models"
	"github.com/officiallysidsingh/ecom-server/internal/store"
)

const userContextKey models.ContextKey = "user"

type OrderService interface {
	GetAll(ctx context.Context) ([]models.Order, error)
	GetByID(ctx context.Context, orderID string) (*models.Order, error)
	// Create(ctx context.Context, order *models.Order) (string, error)
	// PutUpdate(ctx context.Context, order *models.Order, orderID string) error
	// PatchUpdate(ctx context.Context, order *models.Order, orderID string) error
	// Delete(ctx context.Context, orderID string) error
}

type orderService struct {
	store store.OrderStore
}

func NewOrderService(store store.OrderStore) OrderService {
	return &orderService{
		store: store,
	}
}

func (s *orderService) GetAll(ctx context.Context) ([]models.Order, error) {
	// Retrieve user from context
	user, ok := ctx.Value(userContextKey).(models.Claims)

	if !ok {
		log.Println("User not found in context")
		return nil, errors.New("user not found in context")
	}

	// Extract user_id from Claims
	userID := user.UserID

	return s.store.GetAllFromDB(ctx, userID)
}

func (s *orderService) GetByID(ctx context.Context, orderID string) (*models.Order, error) {
	// Retrieve user from context
	user, ok := ctx.Value(userContextKey).(models.Claims)
	if !ok {
		log.Println("User not found in context")
		return nil, fmt.Errorf("user not found in context")
	}

	// Extract user_id from Claims
	userID := user.UserID

	return s.store.GetByIDFromDB(ctx, orderID, userID)
}
