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

func (server *Server) dockerHostList(ctx *gin.Context) {
	hostconfigs, err := server.config.GetDockerHosts()

	if err != nil {
		logger.Log.Error("dockerHostList error.. [%v]", err)
		ctx.JSON(http.StatusInternalServerError, ErrorResponse(err.Error()))
		return
	}

	response := ToContainerHostResponse(hostconfigs)
	ctx.JSON(http.StatusOK, SuccessResponse(response))
}

func (server *Server) dockerPs(ctx *gin.Context) {
	containers, err := server.service.ContainerList(ctx)
	if err != nil {
		logger.Log.Error("Service Container list error.. [%v]", err)
		ctx.JSON(http.StatusInternalServerError, ErrorResponse(err.Error()))
		return
	}

	response := ToContainerListResponse(containers)
	ctx.JSON(http.StatusOK, SuccessResponse(response))
}

type requestHostName struct {
	Host string `uri:"host" binding:"required"`
}

func (server *Server) dockerPs2(ctx *gin.Context) {
	var req requestHostName
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, ErrorResponse(err.Error()))
		return
	}

	containers, err := server.service.ContainerList2(ctx, req.Host)
	if err != nil {
		logger.Log.Error("Service Container list error.. [%v]", err)
		ctx.JSON(http.StatusInternalServerError, ErrorResponse(err.Error()))
		return
	}

	response := ToContainerListResponse(containers)
	ctx.JSON(http.StatusOK, SuccessResponse(response))
}

func (server *Server) containerInspect(ctx *gin.Context) {
	var req requestContainerID
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, ErrorResponse(err.Error()))
		return
	}

	inspect, err := server.service.InspectContainer(ctx, req.ID)
	if err != nil {
		logger.Log.Error("Service Inspect container error.. [%v]", err)
		ctx.JSON(http.StatusInternalServerError, ErrorResponse(err.Error()))
		return
	}

	response := ToContainerInspectResponse(inspect)
	ctx.JSON(http.StatusOK, SuccessResponse(response))
}

type requestHost_ID struct {
	Host string `uri:"host" binding:"required"`
	Id   string `uri:"id" binding:"required"`
}

func (server *Server) containerInspect2(ctx *gin.Context) {
	var req requestHost_ID
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, ErrorResponse(err.Error()))
		return
	}

	inspect, err := server.service.InspectContainer2(ctx, req.Id, req.Host)
	if err != nil {
		logger.Log.Error("Service Inspect container error.. [%v]", err)
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	logger.Log.Print(2, "ID: %s", inspect.ID)
	logger.Log.Print(2, "Image: %s", inspect.Image)
	logger.Log.Print(2, "Name: %s", inspect.Name)

	response := ToContainerInspectResponse(inspect)
	ctx.JSON(http.StatusOK, SuccessResponse(response))
}

func (server *Server) statContainer(ctx *gin.Context) {
	var req requestContainerID
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, ErrorResponse(err.Error()))
		return
	}

	rst, err := server.service.ContainerStats(ctx, req.ID, false)
	if err != nil {
		logger.Log.Error("Service statContainer error.. [%v]", err)
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	logger.Log.Print(2, "cpu : %.2f", rst.CPUPercent)
	logger.Log.Print(2, "memU : %f %s", rst.MemoryUsageVal, rst.MemoryUsageUnit)
	logger.Log.Print(2, "memL : %f %s", rst.MemoryLimitVal, rst.MemoryLimitUnit)
	logger.Log.Print(2, "memP : %.2f", rst.MemoryPercent)
	logger.Log.Print(2, "rx : %d", rst.NetworkRx)
	logger.Log.Print(2, "tx : %d", rst.NetworkTx)

	response := ToContainerStatsResponse(rst)
	ctx.JSON(http.StatusOK, SuccessResponse(response))	
}

func (server *Server) statContainer2(ctx *gin.Context) {
	var req requestHost_ID
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, ErrorResponse(err.Error()))
		return
	}

	rst, err := server.service.ContainerStats2(ctx, req.Id, req.Host, false)
	if err != nil {
		logger.Log.Error("Service statContainer error.. [%v]", err)
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	logger.Log.Print(2, "cpu : %.2f", rst.CPUPercent)
	logger.Log.Print(2, "memU : %f %s", rst.MemoryUsageVal, rst.MemoryUsageUnit)
	logger.Log.Print(2, "memL : %f %s", rst.MemoryLimitVal, rst.MemoryLimitUnit)
	logger.Log.Print(2, "memP : %.2f", rst.MemoryPercent)
	logger.Log.Print(2, "rx : %d", rst.NetworkRx)
	logger.Log.Print(2, "tx : %d", rst.NetworkTx)

	response := ToContainerStatsResponse(rst)
	ctx.JSON(http.StatusOK, SuccessResponse(response))
}
