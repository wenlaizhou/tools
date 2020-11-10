package tools

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/wenlaizhou/middleware"
	"time"
)

type PromResult struct {
	Status string `json:"status"`
	Data   struct {
		ResultType string `json:"resultType"`
		Result     []struct {
			Metric map[string]string `json:"metric"`
			Value  []json.RawMessage `json:"value"`
		} `json:"result"`
	} `json:"data"`
}

func PromQuery(service string, express string) (PromResult, error) {
	result := PromResult{}
	code, _, body, err := middleware.Get(fmt.Sprintf(`%s/api/v1/query?query=%s`, service, express))
	if err != nil {
		return result, err
	}
	if code != 200 {
		return result, errors.New(fmt.Sprintf("code is %v", code))
	}
	err = json.Unmarshal(body, &result)
	return result, err
}

//istio_requests_total{destination_service=~"reg-extraction.*"}
func PromQueryRange(service string, express string, step string, begin time.Time, end time.Time) (PromResult, error) {
	result := PromResult{}
	beginStr := begin.Format("2006-01-02T15:04:05.000Z") //time.RFC3339)
	endStr := end.Format("2006-01-02T15:04:05.000Z")
	queryUrl := fmt.Sprintf("%s/api/v1/query_range?query=%s&start=%s&end=%s&step=%s", service, beginStr, endStr, step)
	code, _, body, err := middleware.Get(queryUrl)
	if err != nil {
		return result, err
	}
	if code != 200 {
		return result, errors.New(fmt.Sprintf("code is %v", code))
	}
	err = json.Unmarshal(body, &result)
	return result, err
}
