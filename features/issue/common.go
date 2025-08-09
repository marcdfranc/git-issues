package issue

import "errors"

var (
	errTitleRequired    = errors.New("title is required")
	errBodyRequired     = errors.New("body is required")
	errCreate           = errors.New("could not create issue")
	errUpdate           = errors.New("could not update issue")
	errClose            = errors.New("could not close  issue")
	errNotFound         = errors.New("issue not found")
	errProcessing       = errors.New("error on process response")
	errNumberIsRequered = errors.New("number is required")
)
