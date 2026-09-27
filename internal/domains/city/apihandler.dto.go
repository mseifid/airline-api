package city

type CreateCityRequest struct {
	ProvinceID int64  `json:"provinceId" example:"1"`
	Name       string `json:"name" example:"Tehran"`
}

type UpdateCityRequest struct {
	ProvinceID int64  `json:"provinceId" example:"1"`
	Name       string `json:"name" example:"Tehran"`
}

type CityResponse struct {
	ID         int64  `json:"id" example:"1"`
	ProvinceID int64  `json:"provinceId" example:"1"`
	Name       string `json:"name" example:"Tehran"`
}
