package initialize

import (
	"aurora/internal/tokens"
	"aurora/middlewares"

	"github.com/gin-gonic/gin"
)

func RegisterRouter() *gin.Engine {
	handler := NewHandle(
		checkProxy(),
		readAccessToken(),
	)

	router := gin.Default()
	router.Use(middlewares.Cors)

	router.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "Hello, world!",
		})
	})

	router.GET("/ping", tokens.Ping)

	//router.GET("/puid/get", tokens.GetConfig)
	//router.GET("/puid/set", tokens.SetConfig)
	//router.POST("/auth/session", handler.session)
	//router.POST("/auth/refresh", handler.refresh)
	router.OPTIONS("/v1/chat/completions", optionsHandler)

	authGroup := router.Group("").Use(middlewares.Authorization)
	authGroup.POST("/v1/chat/completions", handler.nightmare)
	//authGroup.GET("/v1/models", handler.engines)
	//authGroup.POST("/backend-api/conversation", handler.chatgptConversation)
	return router
}
