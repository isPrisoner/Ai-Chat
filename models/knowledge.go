package models

import "time"

// Knowledge 知识库条目
type Knowledge struct {
	ID             string    `json:"id" gorm:"primaryKey;type:varchar(255)"`
	Title          string    `json:"title" gorm:"type:varchar(255);not null"`
	Content        string    `json:"content" gorm:"type:text;not null"`
	Vector         string    `json:"vector" gorm:"type:text"`                  // 向量数据，JSON格式存储
	Source         string    `json:"source" gorm:"type:varchar(255)"`          // 数据来源
	Namespace      string    `json:"namespace" gorm:"type:varchar(100);index"` // 知识域命名空间
	EmbeddingModel string    `json:"embedding_model" gorm:"type:varchar(100)"` // 使用的embedding模型版本
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}
