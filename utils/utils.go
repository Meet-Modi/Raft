package utils

import (
	"strconv"
	"time"

	"golang.org/x/exp/rand"
)

func GenerateRaftPeerId(isBootNode bool) string {
	if isBootNode {
		return "boot"
	}

	// TODO: It should check with other node names
	time := uint64(time.Now().UnixMilli())
	rand.Seed(time)
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	result := make([]byte, 10)
	for i := range result {
		result[i] = charset[rand.Intn(len(charset))]
	}

	result = append([]byte(strconv.FormatUint(time, 10)), result...)
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
