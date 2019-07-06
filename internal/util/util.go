package util

import (
	"fmt"
	"log"
	"os"
)

type Util struct{}

// ResolveAddress matches an input address and provides sane defaults
func (u *Util) ResolveAddress(addr string) string {
	if addr == "" {
		if port := os.Getenv("PORT"); port != "" {
			log.Printf("Environment variable PORT=\"%s\"", port)
			return ":" + port
		}
		log.Printf("Environment variable PORT is undefined. Using port :8080 by default")
		return ":8080"

	} else if addr[0] != ':' {
		return fmt.Sprintf(":%s", addr)
	}
	return addr
}
