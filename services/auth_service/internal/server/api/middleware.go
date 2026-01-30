package api

import (
	"auth-service/internal/logger"
	"fmt"
	"net/http"

	"github.com/gdygd/goglib/token"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
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
