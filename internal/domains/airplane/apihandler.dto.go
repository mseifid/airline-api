package airplane

type CreateAirplaneRequest struct {
	Type     string `json:"type" example:"Airbus A320"`
	Capacity int    `json:"capacity" example:"180"`
}

type UpdateAirplaneRequest struct {
	Type     string `json:"type" example:"Airbus A321"`
	Capacity int    `json:"capacity" example:"220"`
}

type AirplaneResponse struct {
	ID       int64  `json:"id" example:"1"`
	Type     string `json:"type" example:"Airbus A320"`
	Capacity int    `json:"capacity" example:"180"`
}