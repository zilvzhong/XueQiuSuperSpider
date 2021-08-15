package qyapi

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

const (
	corpid     = "ww50601218d6d16ba6" //企业ID
	agentId    = "1000002"             //应用ID
	secret     = "Mge644aDSuNe2zWmHrCTiyQJ1EukAgoxtdv4gpCDpVw"             //Secret
)
//企业微信应用消息提醒方法如下
//func SendCardMsg(ToUsers interface{}, title, description, url string) (map[string]interface{}, error) {
func SendCardMsg( description string) (map[string]interface{}, error) {

	qyurl := fmt.Sprintf("https://qyapi.weixin.qq.com/cgi-bin/gettoken?corpid=%s&corpsecret=%s", corpid, secret)
	data, err := httpGetJson(qyurl)
	if err != nil {
		log.Println(err)
		return data, err
	}
	errcode := data["errcode"].(float64)
	if errcode != 0 {
		log.Println(errcode)
		return make(map[string]interface{}), nil
	}
	access_token := data["access_token"]
	//卡片内容，不同类型消息通知以下内容需修改(可参考企业微信开发文档)
	//req := map[string]interface{}{
	//	"touser":  "@all",
	//	"msgtype": "textcard",
	//	"agentid": agentId,
	//	"textcard": map[string]interface{}{
	//		"title":       title,
	//		"description": description,
	//		"url":         url,
	//		"btntext":     btntxt,
	//	},
	//}
	req := map[string]interface{}{
		"touser":  "@all",
		"msgtype": "text",
		"agentid": agentId,
		"text": map[string]interface{}{
			"content" :  "<a href=\"https://xueqiu.com/u/6741634160\">雪球</a>。\n" + description,
		},
	}

	sendurl := fmt.Sprintf("https://qyapi.weixin.qq.com/cgi-bin/message/send?access_token=%s", access_token)
	data, err = httpPostJson(sendurl, req)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return data, nil
}

//封装了http请求的GET和POST方法，其返回值都是map[string]interface{}，方便我们使用。
func httpGetJson(url string) (map[string]interface{}, error) {
	res, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	var data map[string]interface{}
	err = json.NewDecoder(res.Body).Decode(&data)
	if err != nil {
		return nil, err
	}
	return data, nil
}
func httpPostJson(url string, data map[string]interface{}) (map[string]interface{}, error) {
	res, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}
	resp, err := http.Post(url, "application/json", bytes.NewReader(res))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var data2 map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&data)
	if err != nil {
		return nil, err
	}
	return data2, nil
}

