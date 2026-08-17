// Package response defines the general API response envelope shared by all
// handlers and middleware.
package models

// Response is the general API response envelope used by every endpoint.
type Response struct {
	Message string `json:"message"`
	Body    any    `json:"body,omitempty"`
}

// ResponseWithoutBody is the general API response envelope used by every endpoint that not return a body.
type ResponseWithoutBody struct {
	Message string `json:"message"`
}
