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