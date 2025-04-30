package issue

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"git-issues/domain"
	"git-issues/service/client"
)

const (
	strUrlListFormat = "%s/repos/%s/%s/issues"
)

var (
	errRequest      = errors.New("request error")
	errReadResponse = errors.New("read response error")
	errNotFound     = errors.New("not found")
)

type List interface {
	GetIssues() []domain.Issue
}

type ListFeature struct {
	client client.Client
	url    string
}

func NewList(client client.Client, config *domain.Config) *ListFeature {
	f := &ListFeature{
		client: client,
	}
	f.url = fmt.Sprintf(strUrlListFormat, config.APIBaseURL, config.Owner, config.Repo)
	return f
}

func (f *ListFeature) Get() ([]map[string]interface{}, error) {
	response, err := f.client.MakeRequest(http.MethodGet, f.url, nil)
	if err != nil {
		err = errors.Join(err, errRequest)
		return nil, err
	}

	var issues []map[string]interface{}
	err = json.Unmarshal(response, &issues)
	if err != nil {
		err = errors.Join(err, errReadResponse)
		return nil, err
	}

	if len(issues) == 0 {
		return nil, errNotFound
	}

	return issues, nil
}
