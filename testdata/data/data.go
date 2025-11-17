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

/*var tt = domain.Config{
	Token:      "loadToken",
	Owner:      "loadOwner",
	Repo:       "loadRepo",
	Editor:     "loadEditor",
	APIBaseURL: apiBaseUrl,
}
*/
