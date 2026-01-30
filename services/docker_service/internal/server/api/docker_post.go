package api

import (
	"net/http"

	"docker_service/internal/logger"

	"github.com/gin-gonic/gin"
)

func (server *Server) startContainer(ctx *gin.Context) {
	var req requestContainerID
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	err := server.service.StartContainer(ctx, req.ID)
	if err != nil {
		logger.Log.Error("Service startContainer error.. [%v]", err)
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, "")
}

func (server *Server) stopContainer(ctx *gin.Context) {
	var req requestContainerID
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	err := server.service.StopContainer(ctx, req.ID)
	if err != nil {
		logger.Log.Error("Service stopContainer error.. [%v]", err)
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, "")
}

type requestStartStop struct {
	Id   string `json:"id" binding:"required"`
	Host string `json:"host" binding:"required"`
}

func (server *Server) startContainer2(ctx *gin.Context) {
	var req requestStartStop
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	err := server.service.StartContainer2(ctx, req.Id, req.Host)
	if err != nil {
		logger.Log.Error("Service startContainer error.. [%v]", err)
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, "")
}

func (server *Server) stopContainer2(ctx *gin.Context) {
	var req requestStartStop
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	err := server.service.StopContainer2(ctx, req.Id, req.Host)
	if err != nil {
		logger.Log.Error("Service stopContainer error.. [%v]", err)
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, "")
}
