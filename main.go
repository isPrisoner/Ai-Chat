package main

import (
	"AiDemo/router"
	"log"

	initPkg "AiDemo/init"
	"AiDemo/utils"

	"github.com/gin-gonic/gin"
)

func main() {

	// 统一基础初始化（日志、配置、数据库）
	utils.Info("正在进行基础初始化...")
	cleanup, err := initPkg.InitBase()
	if err != nil {
		log.Fatalf("系统初始化失败: %v", err)
	}
	defer cleanup()
	utils.Info("基础初始化完成")

	// 启动 HTTP 服务
	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()

	router.Register(r)

	utils.Info("🚀 服务已启动: http://localhost:8080")

	if err := r.Run(":8080"); err != nil {
		utils.Fatal("服务启动失败: %v", err)
	}
}
