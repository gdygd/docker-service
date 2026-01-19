package api

import (
	"net/http"
	"time"

	"docker_service/internal/logger"

	"github.com/gin-gonic/gin"
)

func (server *Server) testapi(ctx *gin.Context) {
	time.Sleep(time.Microsecond * 3000)

	strdt, err := server.dbHnd.ReadSysdate(ctx)
	if err != nil {
		logger.Log.Error("testapi err..%v", err)
	}
	logger.Log.Print(2, "testapi :%v", strdt)

	ctx.JSON(http.StatusOK, "hello")
}

func (server *Server) heartbeat(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, nil)
}

func (server *Server) terminate(ctx *gin.Context) {
	server.ch_terminate <- true
	logger.Log.Print(2, "Accept terminate command..")
	ctx.JSON(http.StatusOK, nil)
}

func (server *Server) dockerPs(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, "")

	containers, err := server.service.ContainerList(ctx)
	if err != nil {
		logger.Log.Error("Service Container list error.. [%v]", err)
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	for i, c := range containers {
		logger.Log.Print(2, "list[%d], (%v)", i, c)
	}

	ctx.JSON(http.StatusOK, containers)
}

type getInspectRequest struct {
	ID string `uri:"id" binding:"required"`
}

func (server *Server) dockerInspect(ctx *gin.Context) {
	var req getInspectRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	inspect, err := server.service.InspectContainer(ctx, req.ID)
	if err != nil {
		logger.Log.Error("Service Inspect container error.. [%v]", err)
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	logger.Log.Print(2, "ID: %s", inspect.ID)
	logger.Log.Print(2, "Image: %s", inspect.Image)
	logger.Log.Print(2, "Name: %s", inspect.Name)

	ctx.JSON(http.StatusOK, inspect)
}
