package api

type LoginRequestForm struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginResponse struct {
	UUID        string `json:"uuid"`
	Nickname    string `json:"nickname"`
	Email       string `json:"email"`
	AccessToken string `json:"access_token"`
}

type CreateUserRequestForm struct {
	Nickname string `json:"nickname"`
	Email    string `json:"email"`
	Password string `json:"password"`
}
