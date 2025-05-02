package driver

type Driver struct {
	UserID          int64  `json:"user_id"`
	VehicleInfo     string `json:"vehicle_info"`
	CurrentLocation string `json:"current_location"`
	Available       bool   `json:"available"`
}
