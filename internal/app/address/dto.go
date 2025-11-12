package address

// CreateAddressRequest represents the request to create an address
type CreateAddressRequest struct {
	UserID      string  `json:"user_id" validate:"required,uuid" example:"550e8400-e29b-41d4-a716-446655440000"`
	Address     string  `json:"address" validate:"required,max=50" example:"123 Main Street"`
	Floor       string  `json:"floor" validate:"required,max=10" example:"5"`
	UnitNo      string  `json:"unit_no" validate:"required,max=10" example:"A"`
	BlockTower  *string `json:"block_tower" validate:"omitempty,max=25" example:"Tower B"`
	CompanyName *string `json:"company_name" validate:"omitempty,max=25" example:"ABC Corp"`
}

// UserCreateAddressRequest represents the request for a user to create their own address
type UserCreateAddressRequest struct {
	Address     string  `json:"address" validate:"required,max=50" example:"123 Main Street"`
	Floor       string  `json:"floor" validate:"required,max=10" example:"5"`
	UnitNo      string  `json:"unit_no" validate:"required,max=10" example:"A"`
	BlockTower  *string `json:"block_tower" validate:"omitempty,max=25" example:"Tower B"`
	CompanyName *string `json:"company_name" validate:"omitempty,max=25" example:"ABC Corp"`
}

// UpdateAddressRequest represents the request to update an address (admin)
type UpdateAddressRequest struct {
	Address     *string `json:"address" validate:"omitempty,max=50" example:"456 New Avenue"`
	Floor       *string `json:"floor" validate:"omitempty,max=10" example:"7"`
	UnitNo      *string `json:"unit_no" validate:"omitempty,max=10" example:"B"`
	BlockTower  *string `json:"block_tower" validate:"omitempty,max=25" example:"Tower C"`
	CompanyName *string `json:"company_name" validate:"omitempty,max=25" example:"XYZ Ltd"`
}

// UserUpdateAddressRequest represents the request for a user to update their own address
type UserUpdateAddressRequest struct {
	Address     *string `json:"address" validate:"omitempty,max=50" example:"456 New Avenue"`
	Floor       *string `json:"floor" validate:"omitempty,max=10" example:"7"`
	UnitNo      *string `json:"unit_no" validate:"omitempty,max=10" example:"B"`
	BlockTower  *string `json:"block_tower" validate:"omitempty,max=25" example:"Tower C"`
	CompanyName *string `json:"company_name" validate:"omitempty,max=25" example:"XYZ Ltd"`
}

// AddressResponse represents the address response
type AddressResponse struct {
	ID          string  `json:"id" example:"750e8400-e29b-41d4-a716-446655440002"`
	UserID      string  `json:"user_id" example:"550e8400-e29b-41d4-a716-446655440000"`
	Address     string  `json:"address" example:"123 Main Street"`
	Floor       string  `json:"floor" example:"5"`
	UnitNo      string  `json:"unit_no" example:"A"`
	BlockTower  *string `json:"block_tower,omitempty" example:"Tower B"`
	CompanyName *string `json:"company_name,omitempty" example:"ABC Corp"`
	IsDefault   bool    `json:"is_default" example:"true"`
	CreatedAt   string  `json:"created_at" example:"2024-01-01T12:00:00Z"`
	UpdatedAt   string  `json:"updated_at" example:"2024-01-02T15:30:00Z"`
}

// SetDefaultAddressRequest represents the request to set default address
type SetDefaultAddressRequest struct {
	AddressID string `json:"address_id" validate:"required,uuid" example:"750e8400-e29b-41d4-a716-446655440002"`
}

// ListAddressesRequest represents the request to list addresses
type ListAddressesRequest struct {
	Limit  int32 `json:"limit" validate:"omitempty,min=1,max=100" example:"10"`
	Offset int32 `json:"offset" validate:"omitempty,min=0" example:"0"`
}
