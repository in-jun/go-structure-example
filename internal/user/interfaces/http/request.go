package http

type UpdateProfileRequest struct {
	Name string `json:"name"`
}

type UpdatePasswordRequest struct {
	CurrentPassword string `json:"current_password"`
	NewPassword     string `json:"new_password"`
}
