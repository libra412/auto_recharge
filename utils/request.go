package utils

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
)

//form 表单 发送
func SendPostForm(requestUrl string, params url.Values) []byte {
	resp, err := http.PostForm(requestUrl, params)
	if err != nil {
		return nil
	}
	defer resp.Body.Close()
	response, _ := ioutil.ReadAll(resp.Body)
	return response
}

// 发送短信
func SendMessage(telephone, jsonData, templateCode, signName string) bool {
	params := url.Values{}
	params.Add("telephone", telephone)
	params.Add("json", jsonData)
	params.Add("templateCode", templateCode)
	params.Add("signName", signName)
	resBody := SendPostForm("http://supcommon.bdxxk.com/tools/sendMessage", params)
	if resBody == nil {
		return false
	}
	var v struct {
		Code int
	}
	err := json.Unmarshal([]byte(resBody), &v)
	fmt.Println("短信发送结果", v, err)
	return err == nil && v.Code == 0
}
