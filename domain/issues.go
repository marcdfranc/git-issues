package domain

type Issue struct {
	Number string `json:"number"`
	Title  string `json:"title"`
	Body   string `json:"body,omitempty"`
	State  string `json:"state,omitempty"`
}
