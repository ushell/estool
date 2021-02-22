package elasticsearch

import (
	"context"
	"encoding/json"
	"fmt"
	"bytes"
	"time"

	"es2log/internal/util"
)

type Query struct {
	Index string `json:"index"`
	StartDate string `json:"start_date"`
	EndDate string `json:"end_date"`
	Match map[string]interface{} `json:"match"`
}

func Search(q Query) (Result, error) {
	var r Result
	var response map[string]interface{}
	var buf bytes.Buffer

	if q.Index == "" {
		return r, util.NewError(-1, "index is empty")
	}

	t1, err := time.Parse("2006-01-02", q.StartDate)
	//t2, _ := time.Parse("2006-01-02", q.EndDate)

	var (
		index = q.Index
		match = q.Match
		startDate = t1.Unix()
		//endDate = t2.Unix()
	)

	query := map[string]interface{}{
		"query": map[string]interface{}{
			"bool": map[string]interface{}{
				"must": map[string]interface{}{
					"match": match,
				},
				"filter": map[string]interface{}{
					"range": map[string]interface{}{
						"@timestamp": map[string]int64{
							"gt": startDate,
							//"lt": endDate,
						},
					},
				},
			},
		},
		"size": 10,
	}

	if err := json.NewEncoder(&buf).Encode(query); err != nil {
		return r, util.NewError(-1, fmt.Sprintf("Error encoding query: %s", err.Error()))
	}

	fmt.Println("[*]", buf.String())

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

	return r, nil
}