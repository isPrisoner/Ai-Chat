package handlers

import "AiDemo/services"

// SessionHandler 会话处理器
type SessionHandler struct {
	sessionService *services.SessionService
}

// NewSessionHandler 创建新的会话处理器
func NewSessionHandler() *SessionHandler {
	return &SessionHandler{
		sessionService: services.NewSessionService(),
	}
}
