package util

import (
	"fmt"
	"log"
	"os"
)

type envFetcher interface {
	GetEnvPort() string
}

type osEnvFetcher struct{}

func (f *osEnvFetcher) GetEnvPort() string {
	return os.Getenv("PORT")
}

type Util struct {
	envFetcher
}

func New() *Util {
	return &Util{
		envFetcher: &osEnvFetcher{},
	}
}

// ResolveAddress matches an input address and provides sane defaults
func (u *Util) ResolveAddress(addr string) string {
	if addr == "" {
		if port := u.GetEnvPort(); port != "" {
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
