package models

import (
	"time"
)

// Session 会话模型
type Session struct {
	ID        string     `json:"id" gorm:"primaryKey;type:varchar(255)"`
	Name      string     `json:"name" gorm:"type:varchar(255);not null"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at,omitempty" gorm:"index"`
}

// ChatMessage 聊天消息模型
type ChatMessage struct {
	ID        uint       `json:"id" gorm:"primaryKey;autoIncrement"`
	SessionID string     `json:"session_id" gorm:"type:varchar(255);not null;index"`
	Role      string     `json:"role" gorm:"type:varchar(50);not null"`
	Content   string     `json:"content" gorm:"type:text;not null"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at,omitempty" gorm:"index"`
}

// SessionWithMessageCount 包含消息数量的会话信息
type SessionWithMessageCount struct {
	Session
	MessageCount int64 `json:"message_count"`
}

// CreateSessionRequest 创建会话请求
type CreateSessionRequest struct {
	Name string `json:"name" binding:"required"`
}

// UpdateSessionRequest 更新会话请求
type UpdateSessionRequest struct {
	Name string `json:"name" binding:"required"`
}
