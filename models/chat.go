package models

// ChatRequest 聊天请求体
type ChatRequest struct {
	Message   string `json:"message" binding:"required"`
	Role      string `json:"role"`
	SessionID string `json:"session_id"`
}

// ChatResponse 聊天响应体
type ChatResponse struct {
	Reply     string `json:"reply"`
	SessionID string `json:"session_id"`
}
