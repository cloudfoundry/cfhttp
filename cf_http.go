package cf_http

import (
	"net/http"
	"time"
)

var config Config

type Config struct {
	Timeout time.Duration
}

func Initialize(timeout time.Duration) {
	config.Timeout = timeout
}

func NewClient() *http.Client {
	return &http.Client{
		Timeout: config.Timeout,
	}
}
