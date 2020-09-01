package services

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
)

// "https://api.hbdm.com"如果无法访问请使用："https://api.btcgateway.pro"。
var api_url = "https://api.hbdm.com"
var api_url_bak = "https://api.btcgateway.pro"

// 获取行情
func GetKLine(symbol, period string, params ...int) Response {
	p := url.Values{}
	p.Add("symbol", symbol)
	p.Add("period", period)
	if len(params) == 0 {
		panic("参数不全")
	}
	if len(params) == 1 || len(params) == 3 {
		p.Add("size", strconv.Itoa(params[0]))
	}
	if len(params) == 2 {
		p.Add("from", strconv.Itoa(params[0]))
		p.Add("to", strconv.Itoa(params[1]))
	}
	uri := `/market/history/kline`
	resp, err := http.Get(api_url_bak + uri + "?" + p.Encode())
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	respBody, _ := ioutil.ReadAll(resp.Body)
	var response Response
	json.Unmarshal(respBody, &response)
	return response
}

// 返回内容
type Response struct {
	Ch   string
	Data []struct {
		Id     int
		Vol    float64
		Count  float64
		Open   float64
		Close  float64
		Low    float64
		High   float64
		Amount float64
	}
	Status string
	Ts     string
}
