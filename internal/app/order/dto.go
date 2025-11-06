package order

// CreateOrderRequest represents the request to create an order
type CreateOrderRequest struct {
	UserID      string  `json:"user_id" validate:"required,uuid"`
	OrderNumber string  `json:"order_number" validate:"required,min=3,max=50"`
	Status      string  `json:"status" validate:"omitempty,oneof=pending processing completed cancelled"`
	TotalAmount float64 `json:"total_amount" validate:"required,min=0"`
	Notes       string  `json:"notes" validate:"omitempty"`
}

// UpdateOrderRequest represents the request to update an order
type UpdateOrderRequest struct {
	Status      *string  `json:"status" validate:"omitempty,oneof=pending processing completed cancelled"`
	TotalAmount *float64 `json:"total_amount" validate:"omitempty,min=0"`
	Notes       *string  `json:"notes" validate:"omitempty"`
}

// OrderResponse represents the order response
type OrderResponse struct {
	ID          string  `json:"id"`
	UserID      string  `json:"user_id"`
	OrderNumber string  `json:"order_number"`
	Status      string  `json:"status"`
	TotalAmount float64 `json:"total_amount"`
	Notes       string  `json:"notes,omitempty"`
	CreatedAt   string  `json:"created_at"`
	UpdatedAt   string  `json:"updated_at"`
}

// ListOrdersRequest represents the request to list orders
type ListOrdersRequest struct {
	UserID string `json:"user_id" validate:"omitempty,uuid"`
	Limit  int32  `json:"limit" validate:"omitempty,min=1,max=100"`
	Offset int32  `json:"offset" validate:"omitempty,min=0"`
}
