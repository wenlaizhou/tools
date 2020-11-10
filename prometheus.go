package tools

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/wenlaizhou/middleware"
	"regexp"
	"strconv"
	"time"
)

type PromResult struct {
	Status string `json:"status"`
	Data   struct {
		ResultType string `json:"resultType"`
		Result     []struct {
			Metric map[string]string `json:"metric"`
			Value  []json.RawMessage `json:"value"` // 第一个值为时间戳 time.Unix(int64(math.Round(Value[0])), 0) 即可转换为时间
		} `json:"result"`
	} `json:"data"`
}

type PromRangeResult struct {
	Status string `json:"status"`
	Data   struct {
		ResultType string `json:"resultType"`
		Result     []struct {
			Metric map[string]string   `json:"metric"`
			Values [][]json.RawMessage `json:"values"` // 第一个值为时间戳 time.Unix(int64(math.Round(Value[0])), 0) 即可转换为时间
		} `json:"result"`
	} `json:"data"`
}

func RawToStr(val json.RawMessage) string {
	if val == nil || len(val) <= 0 {
		return ""
	}
	res := string(val)
	checker := regexp.MustCompile(`^"(.*?)"$`)
	if !checker.MatchString(res) {
		return res
	}
	return res[1 : len(res)-1]
}

func RawToTime(val json.RawMessage) (time.Time, error) {
	res, err := strconv.ParseInt(RawToStr(val), 10, 0)
	if err != nil {
		return time.Time{}, err
	}
	return time.Unix(res, 0), nil
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
func PromQueryRange(service string, express string, step string, begin time.Time, end time.Time) (PromRangeResult, error) {
	result := PromRangeResult{}
	beginStr := begin.Format("2006-01-02T15:04:05.000Z") //time.RFC3339)
	endStr := end.Format("2006-01-02T15:04:05.000Z")
	queryUrl := fmt.Sprintf("%s/api/v1/query_range?query=%s&start=%s&end=%s&step=%s", service, express, beginStr, endStr, step)
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
