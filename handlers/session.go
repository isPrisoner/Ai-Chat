package handlers

import (
	"AiDemo/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

// GetSessions 获取所有会话
func (h *SessionHandler) GetSessions(c *gin.Context) {
	sessions, err := h.sessionService.GetAllSessions()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "获取会话列表失败: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"sessions": sessions,
	})
}

// CreateSession 创建新会话
func (h *SessionHandler) CreateSession(c *gin.Context) {
	var req models.CreateSessionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "请求参数错误: " + err.Error(),
		})
		return
	}

	session, err := h.sessionService.CreateSession(req.Name)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "创建会话失败: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"session": session,
	})
}

// UpdateSession 更新会话名称
func (h *SessionHandler) UpdateSession(c *gin.Context) {
	sessionID := c.Param("id")
	if sessionID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "会话ID不能为空",
		})
		return
	}

	var req models.UpdateSessionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "请求参数错误: " + err.Error(),
		})
		return
	}

	if err := h.sessionService.UpdateSession(sessionID, req.Name); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "更新会话失败: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "会话更新成功",
	})
}

// DeleteSession 删除会话（软删除）
func (h *SessionHandler) DeleteSession(c *gin.Context) {
	sessionID := c.Param("id")
	if sessionID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "会话ID不能为空",
		})
		return
	}

	if err := h.sessionService.DeleteSession(sessionID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "删除会话失败: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "会话删除成功",
	})
}

// GetSessionMessages 获取会话的所有消息
func (h *SessionHandler) GetSessionMessages(c *gin.Context) {
	sessionID := c.Param("id")
	if sessionID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "会话ID不能为空",
		})
		return
	}

	messages, err := h.sessionService.GetSessionMessages(sessionID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "获取消息失败: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"messages": messages,
	})
}

// GetSession 获取单个会话信息
func (h *SessionHandler) GetSession(c *gin.Context) {
	sessionID := c.Param("id")
	if sessionID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "会话ID不能为空",
		})
		return
	}

	session, err := h.sessionService.GetSession(sessionID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "会话不存在: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"session": session,
	})
}
