package address

// CreateAddressRequest represents the request to create an address
type CreateAddressRequest struct {
	UserID      string  `json:"user_id" validate:"required,uuid"`
	Address     string  `json:"address" validate:"required,max=50"`
	Floor       string  `json:"floor" validate:"required,max=10"`
	UnitNo      string  `json:"unit_no" validate:"required,max=10"`
	BlockTower  *string `json:"block_tower" validate:"omitempty,max=25"`
	CompanyName *string `json:"company_name" validate:"omitempty,max=25"`
}

// CreateAddressForUserRequest represents the request for a user to create their own address
type CreateAddressForUserRequest struct {
	Address     string  `json:"address" validate:"required,max=50"`
	Floor       string  `json:"floor" validate:"required,max=10"`
	UnitNo      string  `json:"unit_no" validate:"required,max=10"`
	BlockTower  *string `json:"block_tower" validate:"omitempty,max=25"`
	CompanyName *string `json:"company_name" validate:"omitempty,max=25"`
}

// UpdateAddressRequest represents the request to update an address
type UpdateAddressRequest struct {
	Address     *string `json:"address" validate:"omitempty,max=50"`
	Floor       *string `json:"floor" validate:"omitempty,max=10"`
	UnitNo      *string `json:"unit_no" validate:"omitempty,max=10"`
	BlockTower  *string `json:"block_tower" validate:"omitempty,max=25"`
	CompanyName *string `json:"company_name" validate:"omitempty,max=25"`
}

// AddressResponse represents the address response
type AddressResponse struct {
	ID          string  `json:"id"`
	UserID      string  `json:"user_id"`
	Address     string  `json:"address"`
	Floor       string  `json:"floor"`
	UnitNo      string  `json:"unit_no"`
	BlockTower  *string `json:"block_tower,omitempty"`
	CompanyName *string `json:"company_name,omitempty"`
	CreatedAt   string  `json:"created_at"`
	UpdatedAt   string  `json:"updated_at"`
}

// SetDefaultAddressRequest represents the request to set default address
type SetDefaultAddressRequest struct {
	AddressID string `json:"address_id" validate:"required,uuid"`
}

// ListAddressesRequest represents the request to list addresses
type ListAddressesRequest struct {
	Limit  int32 `json:"limit" validate:"omitempty,min=1,max=100"`
	Offset int32 `json:"offset" validate:"omitempty,min=0"`
}
