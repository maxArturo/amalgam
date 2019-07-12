package util

import (
	"errors"
	"fmt"
	"log"
	"os"
	"strconv"
)

type envFetcher interface {
	getEnv(key string) string
}

type osEnvFetcher struct{}

func (f *osEnvFetcher) getEnv(key string) string {
	return os.Getenv(key)
}

func (u *Util) GetEnvVarInt(envVar string) (int, error) {
	val := u.envFetcher.getEnv(envVar)

	if val == "" {
		return -1, errors.New("no value found")
	}
	numVal, err := strconv.Atoi(val)
	if err != nil {
		return -1, err
	} else if numVal < 0 {
		return -1, errors.New("Expected a positive int")
	}
	return numVal, nil
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
		if port := u.envFetcher.getEnv("PORT"); port != "" {
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
