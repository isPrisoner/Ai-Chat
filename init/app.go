package init

import (
	"AiDemo/config"
	"fmt"
)

// InitBase 完成应用的基础初始化（日志、配置、数据库）
// 返回一个清理函数，负责在程序退出时释放资源。
func InitBase() (func(), error) {
	// 初始化日志
	if err := InitLog(); err != nil {
		return nil, fmt.Errorf("日志系统初始化失败: %w", err)
	}

	// 如果后续任一初始化失败，需要确保已经初始化的资源被正确清理
	cleanup := func() {
		config.CloseDatabase()
		CloseLog()
	}

	// 加载配置
	if err := config.LoadEnv(); err != nil {
		cleanup()
		return nil, fmt.Errorf("加载配置失败: %w", err)
	}

	// 初始化数据库
	if err := config.InitDatabase(); err != nil {
		cleanup()
		return nil, fmt.Errorf("数据库初始化失败: %w", err)
	}

	return cleanup, nil
}
