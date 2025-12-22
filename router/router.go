package router

import (
	"net/http"

	"AiDemo/handlers"
	"AiDemo/utils"

	"github.com/gin-gonic/gin"
)

func Register(r *gin.Engine) {
	// 静态资源
	r.Static("/web", "./web")
	utils.Info("静态文件路由已配置")

	// 默认首页
	r.GET("/", func(c *gin.Context) {
		c.Redirect(http.StatusFound, "/web/index.html")
	})

	// 聊天接口
	r.POST("/chat", handlers.ChatHandler)
	utils.Info("聊天 API 已注册")

	// 会话管理
	sessionHandler := handlers.NewSessionHandler()

	api := r.Group("/api")
	{
		sessions := api.Group("/sessions")
		{
			sessions.GET("", sessionHandler.GetSessions)
			sessions.POST("", sessionHandler.CreateSession)
			sessions.GET("/:id", sessionHandler.GetSession)
			sessions.PUT("/:id", sessionHandler.UpdateSession)
			sessions.DELETE("/:id", sessionHandler.DeleteSession)
			sessions.GET("/:id/messages", sessionHandler.GetSessionMessages)
		}
	}

	utils.Info("会话管理 API 已注册")
}
