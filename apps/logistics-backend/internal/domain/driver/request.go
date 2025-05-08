package driver

type CreateDriverRequest struct {
	FullName        string `json:"full_name" binding:"required"`
	Email           string `json:"email" binding:"required"`
	VehicleInfo     string `json:"vehicle_info" binding:"required"`
	CurrentLocation string `json:"current_location" binding:"required"`
	Available       string `json:"available" binding:"required"`
}

func (r *CreateDriverRequest) ToDriver() *Driver {
	return &Driver{
		FullName:        r.FullName,
		Email:           r.Email,
		VehicleInfo:     r.VehicleInfo,
		CurrentLocation: r.CurrentLocation,
		Available:       false,
	}
}
