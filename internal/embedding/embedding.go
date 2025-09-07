package embedding

import (
	"hash/fnv"
	"math"
	"strings"
)

// TextToVector: Simple 8D stub embedding for demo purposes.
// Replace with a real embedding model for production.
func TextToVector(text string) []float32 {
	vec := make([]float32, 8)
	text = strings.ToLower(strings.TrimSpace(text))
	words := strings.Fields(text)
	if len(words) == 0 {
		return vec
	}
	vec[0] = float32(math.Min(float64(len(text))/50.0, 1.0))  // Text length
	vec[1] = float32(math.Min(float64(len(words))/10.0, 1.0)) // Word count
	for i, word := range words {
		hash := hashWord(word)
		vec[2+(i%6)] += float32(hash%1000) / 1000.0 // Hash features
	}
	normalize(vec)
	return vec
}

func hashWord(word string) uint32 {
	h := fnv.New32a()
	h.Write([]byte(word))
	return h.Sum32()
}

func normalize(vec []float32) {
	var mag float32
	for _, v := range vec {
		mag += v * v
	}
	mag = float32(math.Sqrt(float64(mag)))
	if mag > 0 {
		for i := range vec {
			vec[i] /= mag
		}
	}
}
