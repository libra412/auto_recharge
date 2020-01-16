package services

import (
	"encoding/json"
	"fmt"
	"github.com/libra412/auto_recharge/utils"
	"io/ioutil"
	"net/http"
	"net/url"
)

var api_url = "https://api.bdxxk.com/api/getApiOrder.json"
var handle_url = "https://api.bdxxk.com/api/handleApiOrder.json"

//
type Response struct {
	Issuccess   bool
	Description string
	Data        []TradeInfo
}

//
type TradeInfo struct {
	TradeId           int
	ProductId         int
	FacePrice         float32
	SellPrice         float32
	TargetAccount     string
	TargetAccountType int
	RechargeMode      int
	RechargeModeName  string
	BuyAmount         int
	TotalSalePrice    float32
	Area              string
	Game              string
	GameName          string
	StockShopId       int
	StockShopName     string
	CustomerIp        string
	CustomerRegion    string
	DealDateTime      string
}

// 获取订单
func GetData(merchantId, key string) []TradeInfo {
	params := url.Values{}
	params.Add("merchantId", merchantId)
	params.Add("count", "100")
	params.Add("sign", signMd5(merchantId+key))
	resp, err := http.PostForm(api_url, params)
	if err != nil {
		fmt.Print(err)
		return nil
	}
	defer resp.Body.Close()
	responseBody, _ := ioutil.ReadAll(resp.Body)
	fmt.Println(string(responseBody))
	if responseBody != nil {
		var v Response
		json.Unmarshal(responseBody, &v)
		return v.Data
	}
	return nil
}

//更新订单
func UpdateData(merchantId, key, tradeId, state, stateInfo string) []byte {
	params := url.Values{}
	params.Add("merchantId", merchantId)
	params.Add("tradeId", tradeId)
	params.Add("status", state)
	params.Add("description", stateInfo)
	params.Add("sign", signMd5(merchantId+tradeId+state+key))
	resp, err := http.PostForm(handle_url, params)
	if err != nil {
		fmt.Print(err)
		return nil
	}
	defer resp.Body.Close()
	responseBody, _ := ioutil.ReadAll(resp.Body)
	return responseBody
}

//
func signMd5(params string) string {
	return utils.Md5([]byte(params))
}
