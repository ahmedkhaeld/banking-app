package account

// CreateAccountRequest represents the payload for creating a new account.
// swagger:model CreateAccountRequest
type CreateAccountRequest struct {
	// Currency of the account. Allowed values: USD, EUR, GBP, JPY, EGP, CAD, AUD.
	// Required: true
	// Example: USD
	Currency string `json:"currency" binding:"required,oneof=USD EUR GBP JPY EGP CAD AUD"`

	// Initial balance of the account.
	// Example: 1000
	Balance *int64 `json:"balance,omitempty"`
}

type CreateAccountResponse struct {
	// ID of the created account.
	// Example: "123e4567-e89b-12d3-a456-426614174000"
	ID string `json:"id"`
	// UserID of the account owner.
	// Example: "123e4567-e89b-12d3-a456-426614174001"
	UserID string `json:"user_id"`
	// Owner of the account.
	// Example: "John Doe"
	Owner string `json:"owner"`
	// Currency of the account.
	// Example: USD
	Currency string `json:"currency"`
	// Balance of the account.
	// Example: 1000
	Balance int64 `json:"balance"`
	// CreatedAt is the timestamp when the account was created.
	// Example: "2023-10-01T12:00:00Z"
	CreatedAt string `json:"created_at"`
}

type AccountBalanceResponse struct {
	// ID of the account.
	ID string `json:"id"`
	// Balance of the account.
	Balance int64 `json:"balance"`
	// Currency of the account.
	Currency string `json:"currency"`
}

// UpdateAccountBalanceRequest represents the payload for updating balance
// swagger:model UpdateAccountBalanceRequest
type UpdateAccountBalanceRequest struct {
	Amount int64 `json:"amount" binding:"required"`
}
