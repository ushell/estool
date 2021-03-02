package elasticsearch

import (
	"os"
	"sort"
	"time"

	"estool/internal/util"
)

func Download(q Query) (Result, error) {
	var (
		r        Result
		filename string
	)

	if q.Index == "" {
		return r, util.NewError(-1, "index is empty")
	}

	filename = util.GetDownloadFilename("")

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
		"size": PageSize,
		"sort": [1]map[string]string{
			0: {
				"_id": "asc",
			},
		},
	}

	for {
		r, err := Fetch(q.Index, query, true)
		if err != nil {
			return r, err
		}

		if r.Hit == 0 {
			break
		}

		sort, err := saveToFile(r.Data, q.Field, filename)
		if err != nil {
			return r, err
		}

		if sort == "" {
			break
		}

		query["search_after"] = []string{sort}
	}

	r.Data = make(map[string]interface{})
	r.Data["filename"] = filename

	return r, nil
}

func saveToFile(r map[string]interface{}, field string, filename string) (string, error) {
	f, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_RDWR, 0644)
	defer f.Close()

	if err != nil {
		return "", err
	}

	hitList := r["hits"].(map[string]interface{})["hits"].([]interface{})
	if len(hitList) == 0 {
		return "", nil
	}

	lastHit := hitList[len(hitList)-1]

	s := lastHit.(map[string]interface{})["sort"]
	sortFlag := s.([]interface{})[0].(string)

	sortMap := make(map[int64]interface{}, len(hitList))
	var sortKey []int64

	if field == "" {
		field = "message"
	}

	for _, hit := range hitList {
		source := hit.(map[string]interface{})["_source"]
		data := source.(map[string]interface{})[field]

		if data == nil {
			return "", util.NewError(1, "`"+field+"` not exists")
		}

		timestamp := source.(map[string]interface{})["@timestamp"]
		t, _ := time.Parse(time.RFC3339Nano, timestamp.(string))

		k := t.UnixNano() / 1000000

		sortMap[k] = data
		sortKey = append(sortKey, k)
	}

	sort.Slice(sortKey, func(i, j int) bool {
		return sortKey[i] < sortKey[j]
	})

	for _, k := range sortKey {
		data := sortMap[k]
		f.WriteString(data.(string) + "\n")
	}

	return sortFlag, nil
}
