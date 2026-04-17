package dto

type VoteRequestDto struct {
	DeviceId   string  `json:"device_id"`
	PhoneModel string  `json:"phone_model"`
	Browser    string  `json:"browser"`
	Breakfast  *int16  `json:"breakfast"`
	Lunch      *int16  `json:"lunch"`
	Dinner     *int16  `json:"dinner"`
	Latitude   *string `json:"geo_latitude"`
	Longitude  *string `json:"geo_longitude"`
}

type VoteResponseDto struct {
	Breakfast *int16 `json:"breakfast"`
	Lunch     *int16 `json:"lunch"`
	Dinner    *int16 `json:"dinner"`
}
