package services

import (
	"AiDemo/config"
	"AiDemo/models"
	"AiDemo/utils"
	"errors"
	"time"

	"gorm.io/gorm"
)

// SessionService 会话服务
type SessionService struct{}

// NewSessionService 创建新的会话服务实例
func NewSessionService() *SessionService {
	return &SessionService{}
}

// CreateSession 创建新会话
func (s *SessionService) CreateSession(name string) (*models.Session, error) {
	session := &models.Session{
		ID:        utils.GenerateSessionID(),
		Name:      name,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := config.DB.Create(session).Error; err != nil {
		return nil, err
	}

	return session, nil
}

// GetSession 获取会话信息
func (s *SessionService) GetSession(sessionID string) (*models.Session, error) {
	var session models.Session
	if err := config.DB.Where("id = ? AND deleted_at IS NULL", sessionID).First(&session).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("会话不存在")
		}
		return nil, err
	}
	return &session, nil
}

// GetAllSessions 获取所有会话（包含消息数量）
func (s *SessionService) GetAllSessions() ([]models.SessionWithMessageCount, error) {
	var sessions []models.SessionWithMessageCount

	err := config.DB.Table("sessions").
		Select("sessions.*, COUNT(chat_messages.id) as message_count").
		Joins("LEFT JOIN chat_messages ON sessions.id = chat_messages.session_id AND chat_messages.deleted_at IS NULL").
		Where("sessions.deleted_at IS NULL").
		Group("sessions.id").
		Order("sessions.updated_at DESC").
		Find(&sessions).Error

	if err != nil {
		return nil, err
	}

	return sessions, nil
}

// UpdateSession 更新会话名称
func (s *SessionService) UpdateSession(sessionID string, name string) error {
	return config.DB.Model(&models.Session{}).
		Where("id = ? AND deleted_at IS NULL", sessionID).
		Updates(map[string]interface{}{
			"name":       name,
			"updated_at": time.Now(),
		}).Error
}

// DeleteSession 软删除会话
func (s *SessionService) DeleteSession(sessionID string) error {
	now := time.Now()
	return config.DB.Model(&models.Session{}).
		Where("id = ? AND deleted_at IS NULL", sessionID).
		Update("deleted_at", now).Error
}

// GetSessionMessages 获取会话的所有消息
func (s *SessionService) GetSessionMessages(sessionID string) ([]models.ChatMessage, error) {
	var messages []models.ChatMessage
	err := config.DB.Where("session_id = ? AND deleted_at IS NULL", sessionID).
		Order("created_at ASC").
		Find(&messages).Error
	return messages, err
}

// GetHistory 获取会话历史（兼容旧接口）
func (s *SessionService) GetHistory(sessionID string) []models.Message {
	messages, err := s.GetSessionMessages(sessionID)
	if err != nil {
		return []models.Message{}
	}

	// 转换为旧的消息格式
	var result []models.Message
	for _, msg := range messages {
		result = append(result, models.Message{
			Role:    msg.Role,
			Content: msg.Content,
		})
	}

	return result
}

// AddMessage 添加消息到会话
func (s *SessionService) AddMessage(sessionID string, role, content string) error {
	message := &models.ChatMessage{
		SessionID: sessionID,
		Role:      role,
		Content:   content,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := config.DB.Create(message).Error; err != nil {
		return err
	}

	// 更新会话的更新时间
	return config.DB.Model(&models.Session{}).
		Where("id = ?", sessionID).
		Update("updated_at", time.Now()).Error
}
