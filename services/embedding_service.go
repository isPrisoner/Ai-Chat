package services

// EmbeddingModelVersion 当前使用的embedding模型版本
const EmbeddingModelVersion = "mock-v1"

// EmbedText 将文本转换为向量
func EmbedText(text string) ([]float64, error) {
	// TODO: 接入真实的embedding服务
	return []float64{0.1, 0.2, 0.3}, nil
}

// EmbedTextBatch 批量向量化，便于大文本批处理
func EmbedTextBatch(texts []string) ([][]float64, error) {
	vectors := make([][]float64, 0, len(texts))
	for _, t := range texts {
		vec, err := EmbedText(t)
		if err != nil {
			return nil, err
		}
		vectors = append(vectors, vec)
	}
	return vectors, nil
}

// GetEmbeddingModelVersion 获取当前 Embedding 模型版本
func GetEmbeddingModelVersion() string {
	return EmbeddingModelVersion
}
