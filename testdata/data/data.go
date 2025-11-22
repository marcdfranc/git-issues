package data

import "git-issues/domain"

var (
	DefaultConfig = domain.Config{
		Token:      "myToken",
		Owner:      "myOwner",
		Repo:       "myRepo",
		Editor:     "myEditor",
		APIBaseURL: domain.ApiBaseUrl,
	}
)
