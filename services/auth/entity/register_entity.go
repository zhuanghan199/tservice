package entity
type RegisterRequest struct {
	Name     string `json:"name"`
	Pwd      string
	Email    string `json:"email"`
	Phone    string `json:"phone"`
}
