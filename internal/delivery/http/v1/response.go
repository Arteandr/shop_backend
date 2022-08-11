package v1

type ErrorResponse struct {
	Error string `json:"error"`
}

type UserResponse struct {
	Id    int    `json:"id"`
	Email string `json:"email"`
}

type IdResponse struct {
	Id int `json:"id"`
}
