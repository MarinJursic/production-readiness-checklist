package safeclient

import (
	"net/http"
	"time"
)

func fetch() (*http.Response, error) {
	client := http.Client{Timeout: 10 * time.Second}
	return client.Get("https://example.invalid")
}
