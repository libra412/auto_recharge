package main

import (
	"fmt"
	"github.com/andlabs/ui"
	_ "github.com/andlabs/ui/winmanifest"
	"github.com/astaxie/beego/toolbox"
	"github.com/sirupsen/logrus"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"

	"github.com/libra412/auto_recharge/services"
)

var log = logrus.New()

var merchantId, key, productId string

var isCanRequest = true
var isQB = 0

var (
	keyBack     = "input keyevent 4"
	inputSecret = "input text 000000"
)

//
func main() {
	log.Out = os.Stdout //日志标准输出
	file, err := os.OpenFile("logs/wx.log", os.O_CREATE|os.O_WRONLY, 1)
	if err == nil {
		log.Out = file
	} else {
		log.Info("failed to log to file")
	}
	tk := toolbox.NewTask("tk", "0/5 * * * * ?", f)
	toolbox.AddTask("tk", tk)
	ui.Main(setUp)
}

//
func setUp() {
	// 初始化窗口
	mainwin := ui.NewWindow("挂机软件V0.0.1.beta", 640, 480, false)
	mainwin.OnClosing(func(*ui.Window) bool {
		ui.Quit()
		return true
	})
	ui.OnShouldQuit(func() bool {
		mainwin.Destroy()
		return true
	})

	//tab
	tab := ui.NewTab()
	tab.Append("QB/DNF", makeControl(mainwin))
	tab.SetMargined(0, true)
	// tab.Append("第二页", newBox())
	// tab.SetMargined(1, true)
	// 设置tab页
	mainwin.SetChild(tab)
	mainwin.SetMargined(true)

	// 最后显示
	mainwin.Show()
}

//
func makeControl(w *ui.Window) ui.Control {
	hbox := ui.NewHorizontalBox()
	hbox.SetPadded(true)

	entryForm := ui.NewForm()
	entryForm.SetPadded(true)
	hbox.Append(entryForm, false)
	// 设置发送请求按钮
	choosed := ui.NewRadioButtons()
	choosed.Append("QB")
	choosed.Append("DNF")
	entryForm.Append("类型", choosed, false)
	secretInput := ui.NewEntry()
	entryForm.Append("输入支付密码", secretInput, false)
	isRequest := false
	requestButton := ui.NewButton("开始接单")
	entryForm.Append("", requestButton, false)
	//
	supInput := ui.NewEntry()
	entryForm.Append("SUP商户号", supInput, false)
	supSecretInput := ui.NewEntry()
	entryForm.Append("SUP商户密钥", supSecretInput, false)
	productIdInput := ui.NewEntry()
	entryForm.Append("商品ID", productIdInput, false)

	testButton := ui.NewButton("手动充值")
	entryForm.Append("", testButton, false)
	accountInput := ui.NewEntry()
	entryForm.Append("输入充值账号", accountInput, false)
	moneyInput := ui.NewEntry()
	entryForm.Append("输入充值金额", moneyInput, false)
	orderIdInput := ui.NewEntry()
	entryForm.Append("输入订单号", orderIdInput, false)

	//
	testButton.OnClicked(func(*ui.Button) {
		account := accountInput.Text()
		wxSecret := secretInput.Text()
		money := moneyInput.Text()
		orderId := orderIdInput.Text()
		isQB = choosed.Selected()
		if len(account) == 0 {
			ui.MsgBoxError(w, "错误提示", "账号不能为空")
			return
		}
		if len(wxSecret) == 0 {
			ui.MsgBoxError(w, "错误提示", "密码不能为空")
			return
		} else {
			inputSecret = "input text " + wxSecret
		}
		if len(money) == 0 {
			ui.MsgBoxError(w, "错误提示", "金额不能为空")
		}
		if len(orderId) == 0 {
			orderId = "123123123"
		}
		if isQB == 0 {
			go rechargeQB2(account, money, orderId)
		} else {
			go rechargeDNF2(account, money+"00", orderId)
		}

	})
	//
	requestButton.OnClicked(func(*ui.Button) {
		if isRequest {
			isRequest = false
			isCanRequest = false
			toolbox.StopTask()
			requestButton.SetText("开始接单")
		} else {
			isRequest = true
			isCanRequest = true
			merchantId = supInput.Text()
			key = supSecretInput.Text()
			productId = productIdInput.Text()
			wxSecret := secretInput.Text()
			isQB = choosed.Selected()
			if len(wxSecret) == 0 {
				ui.MsgBoxError(w, "错误提示", "微信密码不能为空")
				return
			} else {
				inputSecret = "input text " + wxSecret
			}
			//
			toolbox.StartTask()
			requestButton.SetText("停止接单")
		}
	})
	return hbox
}

// 自动接单
func f() error {
	if isCanRequest {
		data := services.GetData(merchantId, key)
		fmt.Println("请求数据", data)
		count := len(data)
		for i := 0; i < count; i++ {
			if strconv.Itoa(data[i].ProductId) == productId && isCanRequest {
				isCanRequest = false
				services.UpdateData(merchantId, key, strconv.Itoa(data[i].TradeId), "1", "处理中")

				code, stateInfo := func() (string, string) {
					if isQB == 0 {
						return rechargeQB2(data[i].TargetAccount, strconv.Itoa(data[i].BuyAmount), strconv.Itoa(data[i].TradeId))
					} else {
						return rechargeDNF2(data[i].TargetAccount, strconv.Itoa(data[i].BuyAmount*100), strconv.Itoa(data[i].TradeId))
					}
				}()
				if code == "200" {
					services.UpdateData(merchantId, key, strconv.Itoa(data[i].TradeId), "2", stateInfo)
				} else {
					services.UpdateData(merchantId, key, strconv.Itoa(data[i].TradeId), "3", stateInfo)
				}
				time.Sleep(time.Second)
				isCanRequest = true
				break
			}
		}
	} else {
		fmt.Println("有订单在处理中")
	}
	return nil
}

//
var commandsQB = []string{"input roll 100 100", "input tap 1000 920", "input tap 1000 920", "input tap 1000 920", "inputAccount", "input tap 1000 1200",
	"input tap 1000 1200", "input keyevent --longpress 67", " input keyevent --longpress 67", "inputMoney", "input tap 900 950"}

var commandsOppoQB = []string{"input tap 680 700", "input tap 600 920", "input tap 600 1100", "input tap 168 650"}

//
func rechargeQB2(account, money, orderId string) (string, string) {
	fmt.Println("充值QB2")
	begin := time.Now().Unix()
	fileName := orderId + ".png"
	count := len(commandsQB)
	for i := 0; i < count; i++ {
		cmd := commandsQB[i]
		if cmd == "inputAccount" {
			cmd = "input text " + account
		} else if cmd == "inputMoney" {
			cmd = "input text " + money
		}
		execCommandRun(cmd)
		time.Sleep(time.Second)
	}
	time.Sleep(3 * time.Second)
	fmt.Println("开始支付")
	execCommandRun(inputSecret)
	//
	time.Sleep(2 * time.Second)
	return checkResult(fileName, account, orderId, money, begin)
}

//
var commandsDNF = []string{"input roll 100 100", "input tap 80 1800", "input tap 1000 600", "input tap 1000 600", "inputAccount", "input tap 1000 400",
	"input tap 100 1400", "inputMoney", "input tap 900 1100"}

//
func rechargeDNF2(account, money, orderId string) (string, string) {
	fmt.Println("充值DNF2")
	begin := time.Now().Unix()
	fileName := orderId + ".png"
	count := len(commandsDNF)
	for i := 0; i < count; i++ {
		cmd := commandsDNF[i]
		if cmd == "inputAccount" {
			cmd = "input text " + account
		} else if cmd == "inputMoney" {
			cmd = "input text " + money
		}
		execCommandRun(cmd)
		time.Sleep(time.Second)
	}
	time.Sleep(3 * time.Second)
	fmt.Println("开始支付")
	execCommandRun(inputSecret)
	//
	time.Sleep(2 * time.Second)
	return checkResult(fileName, account, orderId, money, begin)

}

/*
func rechargeDNF(account, money, orderId string) (string, string) {
	fmt.Println("充值DNF")
	begin := time.Now().Unix()
	//
	clickEmpty := "input tap 168 852"
	clickDnf := "input tap 10 1400"
	oneAccount := "input tap 680 400 "
	inputAccount := "input text " + account
	clickMoney := "input tap 168 952"
	inputMoney := "input text " + money
	clickPay := "input tap 650 1452"
	fileName := orderId + ".png"
	//
	execCommandRun(clickDnf)
	time.Sleep(time.Second)
	execCommandRun(oneAccount)
	execCommandRun(oneAccount)
	execCommandRun(inputAccount)
	execCommandRun(clickEmpty)
	execCommandRun(clickMoney)
	execCommandRun(inputMoney)
	execCommandRun(clickEmpty)
	execCommandRun(clickPay)
	time.Sleep(3 * time.Second)
	fmt.Println("开始支付")
	execCommandRun(inputSecret)
	//
	time.Sleep(2 * time.Second)
	return checkResult(fileName, account, orderId, money, begin)
}*/

//
func checkResult(fileName, account, orderId, money string, begin int64) (string, string) {
	screencapImage := "screencap -p /data/local/tmp/" + fileName
	copyImage := "/data/local/tmp/" + fileName
	desImage := "./images/" + fileName
	fmt.Println("截图认证")
	execCommandRun(screencapImage)
	execCommand(copyImage, desImage)
	res, err := execCommandSuccess(desImage)
	if err == nil {
		if strings.Contains(res, "支 付 成 功") {
			execCommandRun(keyBack)
			if isQB == 0 {
				execCommandRun("input tap 500 902")
			} else {
				execCommandRun(keyBack)
			}
			log.Info("订单号:", orderId, "，耗时：", time.Now().Unix()-begin, "，充值账号：", account, "，充值金额：", money, "，处理成功")
			return "200", "充值成功"
		} else if strings.Contains(res, "温 馨 提 示") {
			log.Error("订单号:", orderId, "，耗时：", time.Now().Unix()-begin, "，充值账号：", account, "，充值金额：", money, "，错误信息：", "账号错误")
			execCommandRun(keyBack)
			return "202", "账号错误"
		} else if strings.Contains(res, "你 已 在 当 前 商 户 支 付 过 一 笔 相 同 金 额") { // 继续执行
			execCommandRun("input tap 600 900")
			time.Sleep(3 * time.Second)
			fmt.Println("开始支付")
			execCommandRun(inputSecret)
			//
			time.Sleep(2 * time.Second)
			return checkResult(fileName, account, orderId, money, begin)
		} else {
			log.Error("订单号:", orderId, "，耗时：", time.Now().Unix()-begin, "，充值账号：", account, "，充值金额：", money, "，错误信息：", res)
			execCommandRun(keyBack)
			return "203", "未知错误"
		}
	} else {
		log.Error("订单号:", orderId, "，耗时：", time.Now().Unix()-begin, "，充值账号：", account, "，充值金额：", money, "，错误信息：", err)
		return "999", "系统异常"
	}
}

// shell执行命令
func execCommandRun(cmd string) error {
	c := exec.Command("adb", "shell", cmd)
	err := c.Run()
	fmt.Println(err)
	return err
}

// 执行命令
func execCommand(org, des string) error {
	c := exec.Command("adb", "pull", org, des)
	err := c.Start()
	fmt.Println("download image", err)
	return err
}

// 识别图片
func execCommandSuccess(fileName string) (string, error) {
	c := exec.Command("tesseract", fileName, "stdout", "-l", "chi_sim")
	res, err := c.Output()
	return string(res), err
}
