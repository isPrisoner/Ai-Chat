package config

import (
	"AiDemo/models"
	"AiDemo/utils"
	"os"
	"path/filepath"

	"github.com/glebarez/sqlite" // ✅ 纯 Go sqlite 驱动
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

// InitDatabase 初始化数据库
func InitDatabase() error {
	// 确保数据目录存在
	dataDir := "./data"
	if err := os.MkdirAll(dataDir, 0755); err != nil {
		return err
	}

	dbPath := filepath.Join(dataDir, "chat.db")
	// 使用 pure-go sqlite 驱动，不依赖 cgo
	dsn := dbPath + "?_pragma=foreign_keys(1)&_pragma=journal_mode(WAL)"

	var err error
	DB, err = gorm.Open(sqlite.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		return err
	}

	// 自动迁移数据库表（会话、消息、知识库）
	if err := DB.AutoMigrate(
		&models.Session{},
		&models.ChatMessage{},
		&models.Knowledge{},
	); err != nil {
		return err
	}

	utils.Info("数据库初始化完成")
	return nil
}

// CloseDatabase 关闭数据库连接
func CloseDatabase() {
	if DB != nil {
		sqlDB, err := DB.DB()
		if err == nil {
			err := sqlDB.Close()
			if err != nil {
				return
			}
		}
	}
}
