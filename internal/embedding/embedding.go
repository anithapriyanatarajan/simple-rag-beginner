package embedding

import (
	"hash/fnv"
	"math"
	"strings"
)

// Improved text-to-vector function: creates more meaningful 8D vectors
// Uses word-based features and hashing for better semantic representation
func TextToVector(text string) []float32 {
	vec := make([]float32, 8)

	// Normalize text
	text = strings.ToLower(strings.TrimSpace(text))
	words := strings.Fields(text)

	if len(words) == 0 {
		return vec
	}

	// Feature 0: Text length (normalized)
	vec[0] = float32(math.Min(float64(len(text))/50.0, 1.0))

	// Feature 1: Word count (normalized)
	vec[1] = float32(math.Min(float64(len(words))/10.0, 1.0))

	// Feature 2-7: Word-based features using hashing
	for i, word := range words {
		if len(word) > 0 {
			hash := hashWord(word)
			// Distribute word hashes across remaining dimensions
			vec[2+(i%6)] += float32(hash%1000) / 1000.0
		}
	}

	// Normalize vector to unit length for cosine similarity
	normalize(vec)

	return vec
}

// hashWord creates a hash for a word
func hashWord(word string) uint32 {
	h := fnv.New32a()
	h.Write([]byte(word))
	return h.Sum32()
}

// normalize converts vector to unit length
func normalize(vec []float32) {
	var magnitude float32
	for _, v := range vec {
		magnitude += v * v
	}
	magnitude = float32(math.Sqrt(float64(magnitude)))

	if magnitude > 0 {
		for i := range vec {
			vec[i] /= magnitude
		}
	}
}
