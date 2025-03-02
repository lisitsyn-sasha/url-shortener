package delete

import "url-shortener/internal/lib/api/response"

type Request struct {
	Alias string `json:"alias,omitempty"`
}

type Response struct {
	response.Response
	Alias string `json:"alias,omitempty"`
}
