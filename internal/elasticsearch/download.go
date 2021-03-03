package elasticsearch

import (
	"math/rand"
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
	if q.Field == "" {
		q.Field = "message"
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
				"_id":        "asc",
				"@timestamp": "asc",
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

		s, err := saveToFile(r.Data, q.Field, filename)
		if err != nil {
			return r, err
		}

		if len(s) == 0 {
			break
		}

		query["search_after"] = s
	}

	r.Data = make(map[string]interface{})
	r.Data["filename"] = filename

	return r, nil
}

func saveToFile(r map[string]interface{}, field string, filename string) ([]interface{}, error) {
	var sortFlag []interface{}

	f, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_RDWR, 0644)
	defer f.Close()

	if err != nil {
		return sortFlag, err
	}

	hitList := r["hits"].(map[string]interface{})["hits"].([]interface{})
	if len(hitList) == 0 {
		return sortFlag, nil
	}

	lastHit := hitList[len(hitList)-1]
	sortFlag = lastHit.(map[string]interface{})["sort"].([]interface{})

	sortMap := make(map[int64]interface{}, len(hitList))
	var sortKey []int64

	for i, hit := range hitList {
		source := hit.(map[string]interface{})["_source"]
		data := source.(map[string]interface{})[field]
		if data == nil {
			return sortFlag, util.NewError(1, "`"+field+"` not exists")
		}

		timestamp := source.(map[string]interface{})["@timestamp"]
		t, _ := time.Parse(time.RFC3339Nano, timestamp.(string))

		k := t.UnixNano() / 1000000

		rand.Seed(int64(i))

		if _, ok := sortMap[k]; ok {
			k = k + int64(rand.Intn(1000))
		}

		sortMap[k] = data
		sortKey = append(sortKey, k)
	}

	sort.Slice(sortKey, func(i, j int) bool {
		return sortKey[i] < sortKey[j]
	})

	for _, k := range sortKey {
		f.WriteString(sortMap[k].(string) + "\n")
	}

	return sortFlag, nil
}
