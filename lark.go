package tools

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/wenlaizhou/middleware"
)

var logger = middleware.GetLogger("lark")

func GetToken(appId string, appSecret string) (string, error) {
	code, header, data, err := middleware.PostJsonWithTimeout(30,
		"https://open.feishu.cn/open-apis/auth/v3/tenant_access_token/internal/", map[string]interface{}{
			"app_id":     appId,
			"app_secret": appSecret,
		})
	if err != nil {
		logger.ErrorF("%v, %v, %v", code, header, err.Error())
		return "", err
	}
	var dataMap map[string]interface{}
	err = json.Unmarshal(data, &dataMap)
	if err != nil {
		logger.ErrorF("%v, %v, %v", code, header, err.Error())
		return "", err
	}
	token, hasToken := dataMap["tenant_access_token"]
	if !hasToken {
		return "", errors.New("appId, appSecret错误")
	}
	return fmt.Sprintf("%v", token), nil
}

type groupResult struct {
	Code int                    `json:"code"`
	Data map[string]interface{} `json:"data"`
}

func GetRobotGroup(token string) map[string]string {
	code, _, data, err := middleware.GetFull(30, "https://open.feishu.cn/open-apis/chat/v4/list?page_size=200", map[string]string{
		"Authorization": fmt.Sprintf("Bearer %v", token),
	})
	if err != nil {
		logger.ErrorF("%v %v", code, err.Error())
		return nil
	}
	var dataMap groupResult
	err = json.Unmarshal(data, &dataMap)
	if dataMap.Code != 0 {
		return nil
	}
	groups, hasGroups := dataMap.Data["groups"]
	if !hasGroups {
		return nil
	}
	groupsData, hasData := groups.([]interface{})
	if !hasData {
		return nil
	}
	result := map[string]string{}
	for _, group := range groupsData {
		result[fmt.Sprintf("%v",
			group.(map[string]interface{})["name"])] =

			fmt.Sprintf("%v",
				group.(map[string]interface{})["chat_id"])
	}
	return result
}

func SendLarkMessage(token string, receiver string, message string) {
	msgData := map[string]interface{}{
		"chat_id":  receiver,
		"msg_type": "text",
		"content": map[string]string{
			"text": message,
		},
	}
	dataBytes, _ := json.Marshal(msgData)
	code, _, data, err := middleware.PostFull(30, "https://open.feishu.cn/open-apis/message/v3/send/", map[string]string{
		"Authorization": fmt.Sprintf("Bearer %v", token),
	}, middleware.ApplicationJson, dataBytes)
	if err != nil {
		logger.ErrorF("%v %v", code, err.Error())
		return
	}
	logger.InfoF("发送lark消息: %v", string(data))
}

func SendToAllGroups(token string, receiver string, message string) {
	groups := GetRobotGroup(token)
	if groups == nil {
		return
	}
	for _, group := range groups {
		SendLarkMessage(token, group, message)
	}
}
