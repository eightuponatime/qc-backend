package dto

type VoteMealItemDto struct {
	MealType string  `json:"meal_type"`
	Rating   *int16  `json:"rating"`
	Review   *string `json:"review"`
}

type VoteRequestDto struct {
	DeviceId   string            `json:"device_id"`
	PhoneModel string            `json:"phone_model"`
	Browser    string            `json:"browser"`
	Latitude   *string           `json:"geo_latitude"`
	Longitude  *string           `json:"geo_longitude"`
	Items      []VoteMealItemDto `json:"items"`
}

type VoteMealItemResponseDto struct {
	MealType string  `json:"meal_type"`
	Rating   *int16  `json:"rating"`
	Review   *string `json:"review"`
}

type VoteResponseDto struct {
	Items []VoteMealItemResponseDto `json:"items"`
}
