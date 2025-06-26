package transfer

// DTO for create a transfer
type CreateTransferRequest struct {
	FromAccountID string `json:"from_account_id" binding:"required"`
	ToAccountID   string `json:"to_account_id" binding:"required"`
	Amount        int64  `json:"amount" binding:"required"`
}

// CreateTransferResponse represents the response for creating a transfer.
type CreateTransferResponse struct {
	ID            string `json:"id"`
	FromAccountID string `json:"from_account_id"`
	ToAccountID   string `json:"to_account_id"`
	Amount        int64  `json:"amount"`
	CreatedAt     string `json:"created_at"`
}
