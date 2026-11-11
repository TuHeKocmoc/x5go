package config

import "os"

func HTTPAddr() string {
	if value := os.Getenv("HTTP_ADDR"); value != "" {
		return value
	}
	return ":8080"
}
