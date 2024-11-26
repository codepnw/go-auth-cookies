package models

type UserAuthenReq struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UserRegisterReq struct {
	Name     string `json:"name"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}