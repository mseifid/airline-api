package province

type CreateProvinceRequest struct {
	CountryID int64  `json:"countryId" example:"1"`
	Name      string `json:"name" example:"Tehran"`
}

type UpdateProvinceRequest struct {
	CountryID int64  `json:"countryId" example:"1"`
	Name      string `json:"name" example:"Tehran"`
}

type ProvinceResponse struct {
	ID        int64  `json:"id" example:"1"`
	CountryID int64  `json:"countryId" example:"1"`
	Name      string `json:"name" example:"Tehran"`
}