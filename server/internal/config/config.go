package config

import (
	"log"
	"os"
	"path/filepath"
	"strings"
)

type Config struct {
	Addr      string
	Port      string
	DistDir   string
	AssetDir  string
	IndexFile string
	Favicon   string
	DebugMode bool
}

var AppConfig *Config

func Init() {
	AppConfig = &Config{
		Addr: "0.0.0.0",
		Port: "8080",
	}
	env := os.Getenv("GO_ENV")
	if env == "PROD" {
		if env := os.Getenv("SERVER_ADDR"); env != "" {
			AppConfig.Addr = env
		}
		if env := os.Getenv("SERVER_PORT"); env != "" {
			AppConfig.Port = env
		}
		if env := os.Getenv("DIST_DIR"); env != "" {
			AppConfig.DistDir = env
			AppConfig.AssetDir = filepath.Join(AppConfig.DistDir, "assets")
			AppConfig.IndexFile = filepath.Join(AppConfig.DistDir, "index.html")
			AppConfig.Favicon = filepath.Join(AppConfig.DistDir, "favicon.png")
		}
		if env := os.Getenv("DEBUG_MODE"); env != "" {
			if env == "true" {
				AppConfig.DebugMode = true
			} else {
				AppConfig.DebugMode = false
			}
		}
	} else {
		data, err := os.ReadFile(".env")
		if err != nil {
			log.Fatal(err)
		}
		content := string(data)
		lines := strings.Split(content, "\n")
		for _, line := range lines {
			if len(line) == 0 {
				continue
			}
			line = strings.TrimSuffix(line, "\r")
			env := strings.Split(line, "=")
			if env[0] == "DIST_DIR" {
				AppConfig.DistDir = strings.TrimSuffix(env[1], "\r")
				AppConfig.AssetDir = filepath.Join(AppConfig.DistDir, "assets")
				AppConfig.IndexFile = filepath.Join(AppConfig.DistDir, "index.html")
				AppConfig.Favicon = filepath.Join(AppConfig.DistDir, "favicon.png")
			} else if env[0] == "DEBUG_MODE" {
				if env[1] == "true" {
					AppConfig.DebugMode = true
				} else {
					AppConfig.DebugMode = false
				}
			}
		}
	}
	validate()
}

func validate() {
	if AppConfig.Addr == "" || AppConfig.Port == "" || AppConfig.DistDir == "" || AppConfig.AssetDir == "" ||
		AppConfig.IndexFile == "" || AppConfig.Favicon == "" {
		log.Fatal("Invalid configuration")
	}
}
