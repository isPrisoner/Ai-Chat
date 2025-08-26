package handlers

import (
	"AiDemo/models"
	"AiDemo/services"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// ChatHandler 处理聊天请求
func ChatHandler(c *gin.Context) {
	var requestBody models.ChatRequest

	if err := c.ShouldBindJSON(&requestBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请求参数错误"})
		return
	}

	// 如果没有会话ID，创建新会话
	sessionService := services.NewSessionService()
	var sessionID string
	if requestBody.SessionID == "" {
		// 创建新会话
		session, err := sessionService.CreateSession("新对话 " + time.Now().Format("01-02 15:04"))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "创建会话失败"})
			return
		}
		sessionID = session.ID
	} else {
		sessionID = requestBody.SessionID
	}

	// 保存用户消息到数据库
	if err := sessionService.AddMessage(sessionID, "user", requestBody.Message); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "保存用户消息失败"})
		return
	}

	// 获取系统提示词
	systemPrompt := services.GetSystemPrompt(requestBody.Role)

	// 获取会话历史
	history := sessionService.GetHistory(sessionID)
	if len(history) == 0 {
		// 如果没有历史，添加系统提示词
		history = []models.Message{{Role: "system", Content: systemPrompt}}
	}

	// 构建请求体
	requestData := models.RequestBody{
		Model: "ep-20250811150312-h4mvh", // 使用您的豆包模型ID
		Messages: append(history, models.Message{
			Role:    "user",
			Content: requestBody.Message,
		}),
	}

	// 调用豆包API
	response, err := services.CallDoubao(requestData.Messages)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "调用AI服务失败: " + err.Error()})
		return
	}

	// 保存AI回复到数据库
	if err := sessionService.AddMessage(sessionID, "assistant", response); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "保存AI回复失败"})
		return
	}

	// 返回响应
	c.JSON(http.StatusOK, gin.H{
		"reply":      response,
		"session_id": sessionID,
	})
}
