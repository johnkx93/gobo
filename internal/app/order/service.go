package order

import (
	"context"
	"log/slog"
	"strconv"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/user/coc/internal/audit"
	"github.com/user/coc/internal/db"
	"github.com/user/coc/internal/errors"
)

type Service struct {
	queries      *db.Queries
	auditService *audit.Service
}

func NewService(queries *db.Queries, auditService *audit.Service) *Service {
	return &Service{
		queries:      queries,
		auditService: auditService,
	}
}

// CreateOrder creates a new order
func (s *Service) CreateOrder(ctx context.Context, req CreateOrderRequest) (*OrderResponse, error) {
	userID, err := uuid.Parse(req.UserID)
	if err != nil {
		return nil, errors.Validation("invalid user ID format")
	}

	// Check if user exists
	_, err = s.queries.GetUserByID(ctx, pgtype.UUID{Bytes: userID, Valid: true})
	if err == pgx.ErrNoRows {
		return nil, errors.NotFound("user not found")
	} else if err != nil {
		slog.Error("failed to get user", "user_id", req.UserID, "error", err)
		return nil, errors.Internal("failed to verify user", err)
	}

	// Check if order number already exists
	_, err = s.queries.GetOrderByOrderNumber(ctx, req.OrderNumber)
	if err == nil {
		return nil, errors.AlreadyExists("order with this order number already exists")
	} else if err != pgx.ErrNoRows {
		slog.Error("failed to check existing order", "error", err)
		return nil, errors.Internal("failed to check existing order", err)
	}

	status := req.Status
	if status == "" {
		status = "pending"
	}

	// Create order
	var totalAmount pgtype.Numeric
	totalAmount.ScanScientific(strconv.FormatFloat(req.TotalAmount, 'f', 2, 64))

	order, err := s.queries.CreateOrder(ctx, db.CreateOrderParams{
		UserID:      pgtype.UUID{Bytes: userID, Valid: true},
		OrderNumber: req.OrderNumber,
		Status:      status,
		TotalAmount: totalAmount,
		Notes:       pgtype.Text{String: req.Notes, Valid: req.Notes != ""},
	})
	if err != nil {
		slog.Error("failed to create order", "error", err)
		return nil, errors.Internal("failed to create order", err)
	}

	// Audit log the order creation
	orderID := uuid.UUID(order.ID.Bytes)
	s.auditService.LogCreate(ctx, "orders", orderID, order)

	return toOrderResponse(&order), nil
}

// GetOrder retrieves an order by ID
func (s *Service) GetOrder(ctx context.Context, id string) (*OrderResponse, error) {
	orderID, err := uuid.Parse(id)
	if err != nil {
		return nil, errors.Validation("invalid order ID format")
	}

	order, err := s.queries.GetOrderByID(ctx, pgtype.UUID{Bytes: orderID, Valid: true})
	if err == pgx.ErrNoRows {
		return nil, errors.NotFound("order not found")
	} else if err != nil {
		slog.Error("failed to get order", "id", id, "error", err)
		return nil, errors.Internal("failed to get order", err)
	}

	return toOrderResponse(&order), nil
}

// ListOrders retrieves a list of orders with pagination
func (s *Service) ListOrders(ctx context.Context, limit, offset int32) ([]*OrderResponse, error) {
	if limit <= 0 {
		limit = 10
	}
	if limit > 100 {
		limit = 100
	}

	orders, err := s.queries.ListOrders(ctx, db.ListOrdersParams{
		Limit:  limit,
		Offset: offset,
	})
	if err != nil {
		slog.Error("failed to list orders", "error", err)
		return nil, errors.Internal("failed to list orders", err)
	}

	responses := make([]*OrderResponse, len(orders))
	for i, order := range orders {
		o := order
		responses[i] = toOrderResponse(&o)
	}

	return responses, nil
}

// ListOrdersByUserID retrieves orders for a specific user
func (s *Service) ListOrdersByUserID(ctx context.Context, userID string, limit, offset int32) ([]*OrderResponse, error) {
	uid, err := uuid.Parse(userID)
	if err != nil {
		return nil, errors.Validation("invalid user ID format")
	}

	if limit <= 0 {
		limit = 10
	}
	if limit > 100 {
		limit = 100
	}

	orders, err := s.queries.ListOrdersByUserID(ctx, db.ListOrdersByUserIDParams{
		UserID: pgtype.UUID{Bytes: uid, Valid: true},
		Limit:  limit,
		Offset: offset,
	})
	if err != nil {
		slog.Error("failed to list orders by user", "user_id", userID, "error", err)
		return nil, errors.Internal("failed to list orders", err)
	}

	responses := make([]*OrderResponse, len(orders))
	for i, order := range orders {
		o := order
		responses[i] = toOrderResponse(&o)
	}

	return responses, nil
}

// UpdateOrder updates order information
func (s *Service) UpdateOrder(ctx context.Context, id string, req UpdateOrderRequest) (*OrderResponse, error) {
	orderID, err := uuid.Parse(id)
	if err != nil {
		return nil, errors.Validation("invalid order ID format")
	}

	// Check if order exists
	oldOrder, err := s.queries.GetOrderByID(ctx, pgtype.UUID{Bytes: orderID, Valid: true})
	if err == pgx.ErrNoRows {
		return nil, errors.NotFound("order not found")
	} else if err != nil {
		slog.Error("failed to get order", "id", id, "error", err)
		return nil, errors.Internal("failed to get order", err)
	}

	// Update order
	var totalAmount pgtype.Numeric
	if req.TotalAmount != nil {
		totalAmount.ScanScientific(strconv.FormatFloat(*req.TotalAmount, 'f', 2, 64))
		totalAmount.Valid = true
	}

	order, err := s.queries.UpdateOrder(ctx, db.UpdateOrderParams{
		ID:          pgtype.UUID{Bytes: orderID, Valid: true},
		Status:      pgtype.Text{String: ptrToString(req.Status), Valid: req.Status != nil},
		TotalAmount: totalAmount,
		Notes:       pgtype.Text{String: ptrToString(req.Notes), Valid: req.Notes != nil},
	})
	if err != nil {
		slog.Error("failed to update order", "id", id, "error", err)
		return nil, errors.Internal("failed to update order", err)
	}

	// Audit log the order update
	s.auditService.LogUpdate(ctx, "orders", orderID, oldOrder, order)

	return toOrderResponse(&order), nil
}

// DeleteOrder deletes an order
func (s *Service) DeleteOrder(ctx context.Context, id string) error {
	orderID, err := uuid.Parse(id)
	if err != nil {
		return errors.Validation("invalid order ID format")
	}

	// Check if order exists
	order, err := s.queries.GetOrderByID(ctx, pgtype.UUID{Bytes: orderID, Valid: true})
	if err == pgx.ErrNoRows {
		return errors.NotFound("order not found")
	} else if err != nil {
		slog.Error("failed to get order", "id", id, "error", err)
		return errors.Internal("failed to get order", err)
	}

	err = s.queries.DeleteOrder(ctx, pgtype.UUID{Bytes: orderID, Valid: true})
	if err != nil {
		slog.Error("failed to delete order", "id", id, "error", err)
		return errors.Internal("failed to delete order", err)
	}

	// Audit log the order deletion
	s.auditService.LogDelete(ctx, "orders", orderID, order)

	return nil
}

// Helper functions
func toOrderResponse(order *db.Order) *OrderResponse {
	totalAmount, _ := order.TotalAmount.Float64Value()
	return &OrderResponse{
		ID:          uuid.UUID(order.ID.Bytes).String(),
		UserID:      uuid.UUID(order.UserID.Bytes).String(),
		OrderNumber: order.OrderNumber,
		Status:      order.Status,
		TotalAmount: totalAmount.Float64,
		Notes:       order.Notes.String,
		CreatedAt:   order.CreatedAt.Time.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt:   order.UpdatedAt.Time.Format("2006-01-02T15:04:05Z07:00"),
	}
}

func ptrToString(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}
