package domain

type Config struct {
	Token      string `json:"token"`
	Owner      string `json:"owner"`
	Repo       string `json:"repo"`
	Editor     string `json:"editor,omitempty"`
	APIBaseURL string `json:"api_base_url,omitempty"`
}
