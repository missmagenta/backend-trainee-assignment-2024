package main

import (
	"backend-trainee-assignment-2024/config"
	"log"
	"os"
	"testing"
)

var (
	host     string
	basePath string
)

func TestMain(m *testing.M) {
	cfg := config.RequiredConfig().HTTP
	host = "localhost:"
	dockerHost, exists := os.LookupEnv("HOST")
	if exists {
		host = dockerHost + ":"
	}

	host += cfg.Port
	basePath = "http://" + host

	log.Printf("Integration tests: host %s is available", host)

	code := m.Run()
	os.Exit(code)
}
