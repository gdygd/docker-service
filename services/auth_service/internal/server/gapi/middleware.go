package gapi

import (
	"auth-service/internal/logger"
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gdygd/goglib/token"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func authMiddleware(tokenMaker token.Maker) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		logger.Log.Print(2, "authMiddleware...")

		ctx.Next()
	}
}

func corsMiddleware(origins []string) gin.HandlerFunc {
	fmt.Printf("cors : %v \n", origins)
	return cors.New(cors.Config{
		AllowOrigins: origins,
		// AllowOrigins: []string{"http://localhost:3000", "http://localhost:3001", "http://10.1.0.119:8082", "http://10.1.1.164:8082", "http://theroad.web.com:8082"},
		AllowMethods: []string{
			http.MethodHead,
			http.MethodOptions,
			http.MethodGet,
			http.MethodPost,
			http.MethodPut,
			http.MethodPatch,
			http.MethodDelete,
		},
		AllowHeaders: []string{
			"Content-Type",
			"Authorization",
		},
	})
}

func GrpcServerLogger(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp any, err error) {
	startTime := time.Now()
	result, err := handler(ctx, req)
	duration := time.Since(startTime)

	statusCode := codes.Unknown
	if st, ok := status.FromError(err); ok {
		statusCode = st.Code()
	}

	if err != nil {
		logger.Log.Error("gRPC request Err.. %v", err)
	}

	logger.Log.Print(2, "protocol : grpc, method : %s, status_code : %d status_test : %s duration : %v, received gRPC request",
		info.FullMethod, int(statusCode), statusCode.String(), duration)

	return result, err
}
