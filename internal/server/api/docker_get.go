package api

import (
	"context"
	"net/http"
	"sync"
	"time"

	"docker_service/internal/docker"
	"docker_service/internal/logger"

	"github.com/gin-gonic/gin"
)

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

	response := ToContainerStatsResponse(*rst)
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

	response := ToContainerStatsResponse(*rst)
	ctx.JSON(http.StatusOK, SuccessResponse(response))
}

func (server *Server) statContainer3(ctx *gin.Context) {
	logger.Log.Print(2, "[statContainer3] called - path: %s", ctx.Request.URL.Path)

	var req requestHostName
	if err := ctx.ShouldBindUri(&req); err != nil {
		logger.Log.Error("[statContainer3] bind error: %v", err)
		ctx.JSON(http.StatusBadRequest, ErrorResponse(err.Error()))
		return
	}
	logger.Log.Print(2, "[statContainer3] host: %s", req.Host)

	containers, err := server.service.ContainerList(ctx)
	if err != nil {
		logger.Log.Error("[BR] Service Container list error.. [%v]", err)
		ctx.JSON(http.StatusInternalServerError, ErrorResponse(err.Error()))
		return
	}

	var wg sync.WaitGroup
	var ch_res chan docker.ContainerStats = make(chan docker.ContainerStats, len(containers))

	// timeout을 3초로 줄여서 WriteTimeout(5초) 초과 방지
	child_ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	for _, c := range containers {
		c := c
		wg.Add(1)

		go func() {
			defer wg.Done()
			// child_ctx 사용하여 timeout 적용
			rst, err := server.service.ContainerStats2(child_ctx, req.Host, c.ID, false)
			if err != nil {
				logger.Log.Error("get Containerstats error [%s] [%v]", c.ID, err)
				return
			}

			if rst == nil {
				logger.Log.Error("get Containerstats rst is nil..(%s) ", c.ID)
				return
			}

			rst.ID = c.ID
			rst.Name = c.Name

			select {
			case ch_res <- *rst:
			case <-child_ctx.Done():
			}
		}()
	}

	// 완료 신호 채널
	done := make(chan struct{})
	go func() {
		wg.Wait()
		close(ch_res)
		close(done)
	}()

	// 완료 또는 timeout 대기
	select {
	case <-done:
		logger.Log.Print(2, "[statContainer3] all goroutines completed")
	case <-child_ctx.Done():
		logger.Log.Print(2, "[statContainer3] timeout(3s)")
	}

	var resMap map[string]ContainerStatsResponse = make(map[string]ContainerStatsResponse)
	for r := range ch_res {
		resMap[r.ID] = ToContainerStatsResponse(r)
	}

	ctx.JSON(http.StatusOK, SuccessResponse(resMap))
}
