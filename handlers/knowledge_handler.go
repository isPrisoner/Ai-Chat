package handlers

import (
	"AiDemo/models"
	"AiDemo/services"
	"AiDemo/utils"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

// CreateKnowledgeRequest 知识入库请求体
type CreateKnowledgeRequest struct {
	Title     string `json:"title" binding:"required"`
	Content   string `json:"content" binding:"required"`
	Source    string `json:"source"`
	Namespace string `json:"namespace"`
}

// CreateKnowledgeResponse 知识入库响应体
type CreateKnowledgeResponse struct {
	Knowledges []*models.Knowledge `json:"knowledges"`
	Chunks     int                 `json:"chunks"`
	Message    string              `json:"message"`
}

// CreateKnowledgeHandler 知识入库接口
func CreateKnowledgeHandler(c *gin.Context) {
	var req CreateKnowledgeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请求参数错误: " + err.Error()})
		return
	}

	// 设置默认值
	if req.Source == "" {
		req.Source = "manual"
	}
	if req.Namespace == "" {
		req.Namespace = "default"
	}

	knowledges, err := services.SaveKnowledge(req.Title, req.Content, req.Source, req.Namespace)
	if err != nil {
		utils.Error("知识入库失败: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "知识入库失败: " + err.Error()})
		return
	}

	chunksCount := len(knowledges)
	if chunksCount > 0 {
		utils.Info("知识入库成功: 共 %d 个片段, Title=%s, Namespace=%s", chunksCount, req.Title, req.Namespace)
	}

	message := "知识入库成功，已自动完成向量化"
	if chunksCount > 1 {
		message = fmt.Sprintf("知识入库成功，已自动切分为 %d 个片段并完成向量化", chunksCount)
	}

	c.JSON(http.StatusOK, CreateKnowledgeResponse{
		Knowledges: knowledges,
		Chunks:     chunksCount,
		Message:    message,
	})
}
