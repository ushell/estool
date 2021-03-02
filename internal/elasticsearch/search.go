package elasticsearch

import (
	"estool/internal/util"
)

func Search(q Query) (Result, error) {
	var r Result

	if q.Index == "" {
		return r, util.NewError(-1, "index is empty")
	}

	query := map[string]interface{}{
		"query": map[string]interface{}{
			"bool": map[string]interface{}{
				"must": map[string]interface{}{
					"match": q.Match,
				},
				"filter": map[string]interface{}{
					"range": map[string]interface{}{
						"@timestamp": map[string]string{
							"gte": q.StartDate,
							"lt":  q.EndDate,
						},
					},
				},
			},
		},
		"size": 10,
	}

	return Fetch(q.Index, query, false)
}
