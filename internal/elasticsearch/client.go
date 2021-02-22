package elasticsearch

import (
	es7 "github.com/elastic/go-elasticsearch/v7"
	"net/http"
)

type Client struct {
	Version int
	Instance *es7.Client
}

type Result struct {
	Hit int
	CostTime string
}

var _c Client

func NewClient(url string, username string, password string) (Client, error) {
	tp := http.DefaultTransport.(*http.Transport).Clone()
	tp.TLSClientConfig.InsecureSkipVerify = true

	config := es7.Config{
		Addresses: []string{
			url,
		},
		Username: username,
		Password: password,
		Transport: tp,
	}

	es, err := es7.NewClient(config)

	if err != nil {
		return _c, err
	}

	_c.Instance = es

	return _c, nil

}

func Init(url string, username string, password string) error {
	_, err := NewClient(url, username, password)
	return err
}
