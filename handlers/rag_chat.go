package handlers

import (
	"AiDemo/models"
	"AiDemo/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

// RAGChatRequest RAG 聊天请求体
type RAGChatRequest struct {
	Query     string `json:"query" binding:"required"`
	Mode      string `json:"mode"`
	Namespace string `json:"namespace"`
	TopK      int    `json:"top_k"`
	Debug     bool   `json:"debug"`
}

// RAGChatResponse RAG 聊天响应体
type RAGChatResponse struct {
	Answer    string   `json:"answer"`
	Mode      string   `json:"mode"`
	DocsCount int      `json:"docs_count"`
	Namespace string   `json:"namespace,omitempty"`
	HitDocs   []string `json:"hit_docs,omitempty"`
	Fallback  bool     `json:"fallback,omitempty"`
}

// RAGChatHandler 基于 RAG 的问答接口
func RAGChatHandler(c *gin.Context) {
	var req RAGChatRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请求参数错误: " + err.Error()})
		return
	}

	// 设置默认值
	if req.Mode == "" {
		req.Mode = "rag"
	}
	if req.TopK <= 0 {
		req.TopK = services.DefaultTopK
	}

	if req.Mode == "normal" {
		messages := []models.Message{
			{Role: "user", Content: req.Query},
		}
		answer, err := services.CallDoubao(messages)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "调用 AI 服务失败: " + err.Error()})
			return
		}
		c.JSON(http.StatusOK, RAGChatResponse{
			Answer:    answer,
			Mode:      "normal",
			DocsCount: 0,
		})
		return
	}

	var docs []models.Knowledge
	var err error

	if req.Namespace != "" {
		docs, err = services.RetrieveRelevantDocsByNamespace(req.Query, req.Namespace, req.TopK)
	} else {
		docs, err = services.RetrieveRelevantDocs(req.Query, req.TopK)
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "检索知识库失败: " + err.Error()})
		return
	}

	// 知识库无命中时，退化为普通对话
	if len(docs) == 0 {
		messages := []models.Message{
			{Role: "user", Content: req.Query},
		}
		answer, err := services.CallDoubao(messages)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "调用 AI 服务失败: " + err.Error()})
			return
		}

		c.JSON(http.StatusOK, RAGChatResponse{
			Answer:    answer,
			Mode:      "fallback",
			DocsCount: 0,
			Fallback:  true,
		})
		return
	}

	prompt := services.BuildRAGPrompt(req.Query, docs)

	messages := []models.Message{
		{Role: "system", Content: "你是一个企业级知识库问答助手，请严格根据提供的知识内容回答问题。"},
		{Role: "user", Content: prompt},
	}

	answer, err := services.CallDoubao(messages)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "调用 AI 服务失败: " + err.Error()})
		return
	}

	hitDocs := make([]string, 0, len(docs))
	if req.Debug {
		for _, doc := range docs {
			hitDocs = append(hitDocs, doc.Title)
		}
	}

	c.JSON(http.StatusOK, RAGChatResponse{
		Answer:    answer,
		Mode:      "rag",
		DocsCount: len(docs),
		Namespace: req.Namespace,
		HitDocs:   hitDocs,
		Fallback:  false,
	})
}
