package services

import (
	"AiDemo/models"
	"fmt"
	"strings"
)

// RAGPromptTemplate RAG Prompt 模板配置
type RAGPromptTemplate struct {
	SystemRole     string // 系统角色描述
	KnowledgeIntro string // 知识片段介绍语
	QuestionPrefix string // 问题前缀
	JoinSeparator  string // 知识片段分隔符
}

// DefaultRAGPromptTemplate 默认 RAG Prompt 模板
var DefaultRAGPromptTemplate = RAGPromptTemplate{
	SystemRole:     "你是一个专业助手，请只基于以下知识回答，如果知识中没有相关内容，请明确说明。",
	KnowledgeIntro: "【知识片段 %d - %s】",
	QuestionPrefix: "\n\n问题：",
	JoinSeparator:  "\n\n",
}

// BuildRAGPrompt 根据检索到的文档构建 RAG Prompt
func BuildRAGPrompt(query string, docs []models.Knowledge) string {
	return BuildRAGPromptWithTemplate(query, docs, DefaultRAGPromptTemplate)
}

// BuildRAGPromptWithTemplate 使用自定义模板构建 RAG Prompt
func BuildRAGPromptWithTemplate(query string, docs []models.Knowledge, template RAGPromptTemplate) string {
	if len(docs) == 0 {
		return template.SystemRole + template.QuestionPrefix + query
	}

	var knowledgeParts []string
	for i, d := range docs {
		title := d.Title
		if title == "" {
			title = "未命名"
		}
		intro := fmt.Sprintf(template.KnowledgeIntro, i+1, title)
		knowledgeParts = append(knowledgeParts, intro+"\n"+d.Content)
	}

	knowledgeText := strings.Join(knowledgeParts, template.JoinSeparator)

	return template.SystemRole + "\n\n" + knowledgeText + template.QuestionPrefix + query
}

func joinDocs(docs []models.Knowledge) string {
	var parts []string
	for i, d := range docs {
		parts = append(parts, fmt.Sprintf("【%d】%s\n%s", i+1, d.Title, d.Content))
	}
	return strings.Join(parts, "\n\n")
}
