package utils

import (
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
