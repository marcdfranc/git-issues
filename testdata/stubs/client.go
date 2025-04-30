package stubs

type ClientStub struct {
	MakeRequestFunc func(method, url string, data interface{}) ([]byte, error)
}

func (s *ClientStub) MakeRequest(method, url string, data interface{}) ([]byte, error) {
	if s.MakeRequestFunc != nil {
		return s.MakeRequestFunc(method, url, data)
	}
	return nil, nil
}
