package utils

import (
	"time"

	"golang.org/x/exp/rand"
)

func GenerateRaftPeerId() string {
	// TODO: It should check with other node names
	time := uint64(time.Now().UnixMilli())
	rand.Seed(time)
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	result := make([]byte, 10)
	for i := range result {
		result[i] = charset[rand.Intn(len(charset))]
	}

	return string(result)
}

// RandomKey picks a random key from the given map
func RandomKey(m []string) string {
	if len(m) == 0 {
		return ""
	}

	rand.Seed(uint64(time.Now().UnixMilli()))

	randomIndex := rand.Intn(len(m))
	return m[randomIndex]

}

func Min(a, b int64) int64 {
	if a < b {
		return a
	}
	return b
}
