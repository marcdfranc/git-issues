package domain

const (
	ConfigFile = ".ghissuescli"
	ApiBaseUrl = "https://api.github.com"
)

type Config struct {
	Token      string `json:"token"`
	Owner      string `json:"owner"`
	Repo       string `json:"repo"`
	Editor     string `json:"editor,omitempty"`
	APIBaseURL string `json:"api_base_url,omitempty"`
}
