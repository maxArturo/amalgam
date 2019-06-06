package util

import (
	"log"
	"os"
)

// ResolveAddress matches an input address and provides sane defaults
func ResolveAddress(addr []string) string {
	switch len(addr) {
	case 0:
		if port := os.Getenv("PORT"); port != "" {
			log.Printf("Environment variable PORT=\"%s\"", port)
			return ":" + port
		}
		log.Printf("Environment variable PORT is undefined. Using port :8080 by default")
		return ":8080"
	case 1:
		return addr[0]
	default:
		panic("too many parameters")
	}
}
