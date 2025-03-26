package domain

type Issue struct {
	Title string `json:"title"`
	Body  string `json:"body,omitempty"`
	State string `json:"state,omitempty"`
}