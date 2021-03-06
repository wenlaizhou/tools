package tools

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/wenlaizhou/middleware"
	"math"
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

func PromRawToTime(val json.RawMessage) (time.Time, error) {
	res, err := strconv.ParseFloat(RawToStr(val), 64)
	if err != nil {
		return time.Time{}, err
	}
	return time.Unix(int64(math.Round(res)), 0), nil
}

func RawToInt64(val json.RawMessage) (int64, error) {
	res := RawToStr(val)
	if len(res) <= 0 {
		return -1, errors.New("数据为空")
	}
	return strconv.ParseInt(res, 10, 0)
}

func RawToInt(val json.RawMessage) (int, error) {
	res, err := RawToInt64(val)
	return int(res), err
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
// 标签查询:
// =  : 精确地匹配标签给定的值
// != : 不等于给定的标签值
// =~ : 正则表达匹配给定的标签值
// !~ : 给定的标签值不符合正则表达式
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
