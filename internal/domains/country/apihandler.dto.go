package country

type CreateCountryRequest struct {
	Name string `json:"name" example:"Iran"`
}

type UpdateCountryRequest struct {
	Name string `json:"name" example:"Iran"`
}

type CountryResponse struct {
	ID   int64  `json:"id" example:"1"`
	Name string `json:"name" example:"Iran"`
}