package riskyclient

import web "net/http"

func fetch(client *web.Client) (*web.Response, error) {
	return client.Get("https://example.invalid")
}
