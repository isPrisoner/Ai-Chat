package router

import (
	"net/http"

	"AiDemo/handlers"
	"AiDemo/router/middleware"
	"AiDemo/utils"

	"github.com/gin-gonic/gin"
)

func Register(r *gin.Engine) {
	// 简单限流（可按需调整或关闭）
	r.Use(middleware.RateLimiter(120)) // 每 IP 每分钟 120 次

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

	// RAG 聊天接口（增强版，支持模式区分和多知识域）
	r.POST("/rag/chat", handlers.RAGChatHandler)
	utils.Info("RAG 聊天 API 已注册")

	// 知识入库接口（关键：RAG 从 Demo 到产品的核心接口）
	r.POST("/rag/knowledge", handlers.CreateKnowledgeHandler)
	utils.Info("知识入库 API 已注册")

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
