package api

type createUserRequest struct {
	Username       string `json:"username" binding:"required,alphanum"`
	HashedPassword string `json:"password" binding:"required,min=4"`
	FullName       string `json:"full_name" binding:"required"`
	Email          string `json:"email" binding:"required,email"`
}

type loginUserRequest struct {
	Username string `json:"username" binding:"required,alphanum"`
	Password string `json:"password" binding:"required,min=4"`
}

type renewAccessTokenRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

type logoutUserRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

type requestContainerID struct {
	ID string `uri:"id" binding:"required"`
}

type requestHost_ID struct {
	Host string `uri:"host" binding:"required"`
	Id   string `uri:"id" binding:"required"`
}

type requestHostName struct {
	Host string `uri:"host" binding:"required"`
}
