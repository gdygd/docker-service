package api

import (
	"api-gateway/internal/logger"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/gdygd/goglib/token"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

const (
	authorizationHeaderKey  = "authorization"
	authorizationTypeBearer = "bearer"
	authorizationPayloadKey = "authorization_payload"
)

func authMiddleware(tokenMaker token.Maker) gin.HandlerFunc {
	return func(ctx *gin.Context) {

		path := ctx.Request.URL.Path
		logger.Log.Print(2, " request url : %s", ctx.Request.URL.String())

		if strings.HasPrefix(path, "/auth/login") ||
			strings.HasPrefix(path, "/auth/refresh") ||
			strings.HasPrefix(path, "/public") {
			ctx.Next()
			return
		}

		authorizationHeader := ctx.GetHeader(authorizationHeaderKey)
		if len(authorizationHeader) == 0 {
			err := errors.New("authorization header is not provided")
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse(err))
			return
		}

		fields := strings.Fields(authorizationHeader)
		if len(fields) < 2 {
			err := errors.New("invalid authorization header format")
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse(err))
			return
		}

		authorizationType := strings.ToLower(fields[0])
		if authorizationType != authorizationTypeBearer {
			err := fmt.Errorf("unsupported authorization type %s", authorizationType)
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse(err))
			return
		}

		accessToken := fields[1]
		payload, err := tokenMaker.VerifyToken(accessToken)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse(err))
			return
		}

		ctx.Set(authorizationPayloadKey, payload)
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
