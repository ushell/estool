package elasticsearch

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	es7 "github.com/elastic/go-elasticsearch/v7"

	"estool/internal/util"
)

type Client struct {
	Version  int
	Instance *es7.Client
}

type Result struct {
	Hit      int
	CostTime string
	Data     map[string]interface{}
}

type Query struct {
	Index     string                 `json:"index"`
	StartDate string                 `json:"start_date"`
	EndDate   string                 `json:"end_date"`
	Match     map[string]interface{} `json:"match"`
	Field     string                 `json:"field"`
}

var _c Client

const (
	PageSize = 1000
)

func NewClient(url string, username string, password string) (Client, error) {
	tp := http.DefaultTransport.(*http.Transport).Clone()
	tp.TLSClientConfig.InsecureSkipVerify = true

	config := es7.Config{
		Addresses: []string{
			url,
		},
		Username:  username,
		Password:  password,
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

func Fetch(index string, q map[string]interface{}, isFetch bool) (Result, error) {
	var r Result
	var buf bytes.Buffer
	var response map[string]interface{}

	if err := json.NewEncoder(&buf).Encode(q); err != nil {
		return r, util.NewError(-1, fmt.Sprintf("Error encoding query: %s", err.Error()))
	}

	//fmt.Println("[*] query => ", buf.String())

	res, err := _c.Instance.Search(
		_c.Instance.Search.WithContext(context.Background()),
		_c.Instance.Search.WithIndex(index),
		_c.Instance.Search.WithBody(&buf),
		_c.Instance.Search.WithTrackTotalHits(true),
		_c.Instance.Search.WithPretty(),
	)
	if err != nil {
		return r, util.NewError(-1, fmt.Sprintf("Error getting response: %s", err.Error()))
	}
	defer res.Body.Close()

	if res.IsError() {
		var e map[string]interface{}

		if err := json.NewDecoder(res.Body).Decode(&e); err != nil {
			if res.StatusCode == 401 {
				return r, util.NewError(401, "ES Service need auth !")
			}
			return r, util.NewError(1, fmt.Sprintf("Error parsing the response body: %s", err.Error()))
		} else {
			errorType := e["error"].(map[string]interface{})["type"]
			errorReason := e["error"].(map[string]interface{})["reason"]

			return r, util.NewError(-1, fmt.Sprintf("Error type: %s, reason: %s", errorType, errorReason))
		}
	}

	if err := json.NewDecoder(res.Body).Decode(&response); err != nil {
		return r, util.NewError(-1, fmt.Sprintf("Error parsing the response body: %s", err.Error()))
	}

	hit := int(response["hits"].(map[string]interface{})["total"].(map[string]interface{})["value"].(float64))

	r.Hit = hit
	r.CostTime = fmt.Sprintf("%dms", int(response["took"].(float64)))

	if isFetch == true {
		r.Data = make(map[string]interface{})
		r.Data = response
	}

	return r, nil
}
