package airport

type CreateAirportRequest struct {
	CityID int64  `json:"cityId" example:"1"`
	Code   string `json:"code" example:"IKA"`
	Name   string `json:"name" example:"Imam Khomeini International Airport"`
}

type UpdateAirportRequest struct {
	CityID int64  `json:"cityId" example:"1"`
	Code   string `json:"code" example:"IKA"`
	Name   string `json:"name" example:"Imam Khomeini International Airport"`
}

type AirportResponse struct {
	ID     int64  `json:"id" example:"1"`
	CityID int64  `json:"cityId" example:"1"`
	Code   string `json:"code" example:"IKA"`
	Name   string `json:"name" example:"Imam Khomeini International Airport"`
}