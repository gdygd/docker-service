package api

import (
	"net/http"

	"docker_service/internal/db"
	"docker_service/internal/util"

	"github.com/gin-gonic/gin"
)

func (server *Server) createUser(ctx *gin.Context) {
	var req createUserRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	hashedPassword, err := util.HashPassword(req.HashedPassword)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	arg := db.CreateUserParams{
		Username:       req.Username,
		HashedPassword: hashedPassword,
		FullName:       req.FullName,
		Email:          req.Email,
	}

	user, err := server.service.CreateUser(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	// 응답에서 비밀번호 제외
	createdAt := ""
	if user.CreatedAt.Valid {
		createdAt = user.CreatedAt.Time.Format("2006-01-02 15:04:05")
	}

	resp := struct {
		Username  string `json:"username"`
		FullName  string `json:"full_name"`
		Email     string `json:"email"`
		CreatedAt string `json:"created_at"`
	}{
		Username:  user.Username,
		FullName:  user.FullName,
		Email:     user.Email,
		CreatedAt: createdAt,
	}

	ctx.JSON(http.StatusOK, SuccessResponse(resp))
}