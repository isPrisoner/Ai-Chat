package services

const (
	// DefaultChunkSize 默认文档切分大小（字符数）
	DefaultChunkSize = 500
)

// ChunkText 文档切分，每 chunkSize 字一段
func ChunkText(text string, chunkSize int) []string {
	if chunkSize <= 0 {
		chunkSize = DefaultChunkSize
	}

	if len(text) <= chunkSize {
		return []string{text}
	}

	var chunks []string
	runes := []rune(text)

	for i := 0; i < len(runes); i += chunkSize {
		end := i + chunkSize
		if end > len(runes) {
			end = len(runes)
		}
		chunks = append(chunks, string(runes[i:end]))
	}

	return chunks
}

// ChunkTextWithOverlap 带重叠的文档切分
func ChunkTextWithOverlap(text string, chunkSize, overlapSize int) []string {
	if chunkSize <= 0 {
		chunkSize = DefaultChunkSize
	}
	if overlapSize <= 0 {
		overlapSize = 50
	}

	if len(text) <= chunkSize {
		return []string{text}
	}

	var chunks []string
	runes := []rune(text)

	for i := 0; i < len(runes); {
		end := i + chunkSize
		if end > len(runes) {
			end = len(runes)
		}
		chunks = append(chunks, string(runes[i:end]))
		i = end - overlapSize
		if i >= len(runes) {
			break
		}
	}

	return chunks
}
