package agency

type CreateAgencyRequest struct {
	Name          string `json:"name"`
	WalletBalance *int64 `json:"walletBalance,omitempty"`
}

type CreateAgencyResponse struct {
	ID            int64  `json:"id"`
	Name          string `json:"name"`
	WalletBalance int64  `json:"walletBalance"`
	APIKey        string `json:"apiKey"`
}

type AddBalanceRequest struct {
	Amount int64 `json:"amount" example:"10000000"`
}

type AddBalanceResponse struct {
	WalletBalance int64 `json:"walletBalance"`
}