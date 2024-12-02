package config

import (
	"os"
	"strconv"
)

const DefaultBootNodeURI = "127.0.0.1"
const DefaultDiscoveryProtocol = "gossip" // ["gossip"]
const DefaultIsBootNode = false
const DefaultPort = "3000"

var BootNodeURI = getEnv("BOOT_NODE_URI", DefaultBootNodeURI)
var DiscoveryProtocol = getEnv("DISCOVERY_PROTOCOL", DefaultDiscoveryProtocol)
var IsBootNode, _ = strconv.ParseBool(getEnv("BOOT_NODE", strconv.FormatBool(DefaultIsBootNode)))
var Port = getEnv("PORT", DefaultPort)

var ElectionTimeoutMs = int64(300)
var MaxDiscoveryRetries = 3

// getEnv reads an environment variable or returns a default value if not set
func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}
