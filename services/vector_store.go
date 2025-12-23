package services

import (
	"AiDemo/config"
	"AiDemo/models"
	"encoding/json"
)

// ScoredDoc 检索结果
type ScoredDoc struct {
	Doc   models.Knowledge
	Score float64
}

// VectorStore 向量存储抽象，便于未来替换 Milvus/pgvector/Qdrant
type VectorStore interface {
	Search(queryVec []float64, namespace string, topK int) ([]ScoredDoc, error)
}

// SQLiteVectorStore 基于 SQLite/GORM 的默认实现
type SQLiteVectorStore struct{}

// Search 按余弦相似度检索
func (s *SQLiteVectorStore) Search(queryVec []float64, namespace string, topK int) ([]ScoredDoc, error) {
	var all []models.Knowledge
	db := config.DB
	if namespace != "" {
		db = db.Where("namespace = ?", namespace)
	}
	if err := db.Find(&all).Error; err != nil {
		return nil, err
	}
	if len(all) == 0 {
		return []ScoredDoc{}, nil
	}

	var scored []ScoredDoc
	for _, d := range all {
		if d.Vector == "" {
			continue
		}
		var vec []float64
		if err := json.Unmarshal([]byte(d.Vector), &vec); err != nil {
			continue
		}
		sim := cosineSimilarity(queryVec, vec)
		if sim >= MinSimilarityThreshold {
			scored = append(scored, ScoredDoc{Doc: d, Score: sim})
		}
	}

	if len(scored) == 0 {
		return []ScoredDoc{}, nil
	}

	// 选择排序，按相似度降序取 TopK
	for i := 0; i < len(scored)-1; i++ {
		maxIdx := i
		for j := i + 1; j < len(scored); j++ {
			if scored[j].Score > scored[maxIdx].Score {
				maxIdx = j
			}
		}
		scored[i], scored[maxIdx] = scored[maxIdx], scored[i]
	}

	if topK > 0 && len(scored) > topK {
		scored = scored[:topK]
	}

	return scored, nil
}

var defaultVectorStore VectorStore = &SQLiteVectorStore{}
