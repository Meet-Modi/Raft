package config

import (
	"os"
	"strconv"
)

const DefaultBootNodeURI = "127.0.0.1"
const DefaultDiscoveryProtocol = "gossip" // ["gossip"]
const DefaultIsBootNode = false

var BootNodeURI = getEnv("BOOT_NODE_URI", DefaultBootNodeURI)
var DiscoveryProtocol = getEnv("DISCOVERY_PROTOCOL", DefaultDiscoveryProtocol)
var IsBootNode, _ = strconv.ParseBool(getEnv("BOOT_NODE", strconv.FormatBool(DefaultIsBootNode)))

// getEnv reads an environment variable or returns a default value if not set
func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}
