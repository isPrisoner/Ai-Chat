package main

import (
	"AiDemo/config"
	"AiDemo/handlers"
	initPkg "AiDemo/init"
	"AiDemo/utils"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	// 初始化日志系统
	if err := initPkg.InitLog(); err != nil {
		log.Fatalf("日志系统初始化失败: %v", err)
	}
	defer initPkg.CloseLog()

	// 加载配置
	utils.Info("正在加载配置...")
	err := config.LoadEnv()
	if err != nil {
		utils.Fatal("加载配置失败: %v", err)
		return
	}
	utils.Info("配置加载完成")

	// 初始化数据库
	utils.Info("正在初始化数据库...")
	if err := config.InitDatabase(); err != nil {
		utils.Fatal("数据库初始化失败: %v", err)
		return
	}
	defer config.CloseDatabase()
	utils.Info("数据库初始化完成")

	// 创建 Gin 引擎
	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()

	// 静态文件（前端页面）
	r.Static("/web", "./web")
	utils.Info("静态文件路由已配置")

	// 默认首页跳转
	r.GET("/", func(c *gin.Context) {
		c.Redirect(http.StatusFound, "/web/index.html")
	})

	// 聊天路由
	r.POST("/chat", handlers.ChatHandler)
	utils.Info("聊天API路由已注册")

	// 会话管理路由
	sessionHandler := handlers.NewSessionHandler()
	r.GET("/api/sessions", sessionHandler.GetSessions)
	r.POST("/api/sessions", sessionHandler.CreateSession)
	r.GET("/api/sessions/:id", sessionHandler.GetSession)
	r.PUT("/api/sessions/:id", sessionHandler.UpdateSession)
	r.DELETE("/api/sessions/:id", sessionHandler.DeleteSession)
	r.GET("/api/sessions/:id/messages", sessionHandler.GetSessionMessages)
	utils.Info("会话管理API路由已注册")

	utils.Info("🚀 服务已启动，请在浏览器访问: http://localhost:8080")

	err = r.Run(":8080")
	if err != nil {
		utils.Fatal("服务启动失败: %v", err)
		return
	}
}
