package stubs

import "git-issues/domain"

type ClientStub struct {
	MakeRequestFunc func(method, url string, data *domain.Issue) ([]byte, error)
}

func (s *ClientStub) MakeRequest(method, url string, data *domain.Issue) ([]byte, error) {
	if s.MakeRequestFunc != nil {
		return s.MakeRequestFunc(method, url, data)
	}
	return nil, nil
}
