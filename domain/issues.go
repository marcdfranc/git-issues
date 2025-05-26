package domain

type Issue struct {
	Number int    `json:"number,omitempty"`
	Title  string `json:"title"`
	Body   string `json:"body,omitempty"`
	State  string `json:"state,omitempty"`
}
