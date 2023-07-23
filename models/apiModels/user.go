package apimodels

type User struct {
	Name  string `json:"name"`
	Email string `json:"email"`
	Image string `json:"image"`
	Role  string `json:"role"`
	Id    uint   `json:"id"`
}

type EditUserData struct {
	Name  string `json:"name" validate:"required,max=30"`
	Image string `json:"image"`
}
