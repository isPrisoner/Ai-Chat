package services

func GetSystemPrompt(role string) string {
	switch role {
	case "coder":
		return "你是一个专业的代码专家，擅长各种编程语言和技术栈。请提供清晰、准确的代码建议和解决方案。"
	case "translator":
		return "你是一个专业的翻译官，精通中英文互译。请提供准确、自然的翻译结果。"
	case "pm":
		return "你是一个经验丰富的产品经理，擅长产品设计、需求分析和项目管理。请提供专业的产品建议。"
	case "scholar":
		return "你是一个博学的学术导师，擅长各种学科知识。请提供深入、准确的学术解答。"
	default:
		return "你是一个智能AI助手，能够帮助用户解决各种问题。请提供有用、准确的回答。"
	}
}
