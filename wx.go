package main

import (
	"github.com/andlabs/ui"
	_ "github.com/andlabs/ui/winmanifest"
	"github.com/libra412/auto_recharge/utils"
)

//
func main() {
	ui.Main(setUp)
}

//
func setUp() {
	// 初始化窗口
	mainwin := ui.NewWindow("短信发送0.0.1.beta", 640, 480, false)
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
	tab.Append("短信发送", makeControl(mainwin))
	tab.SetMargined(0, true)
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
	//
	accountInput := ui.NewEntry()
	entryForm.Append("输入充值账号", accountInput, false)
	secretInput := ui.NewEntry()
	entryForm.Append("输入卡号密码", secretInput, false)
	requestButton := ui.NewButton("发送")
	entryForm.Append("", requestButton, false)

	//
	requestButton.OnClicked(func(*ui.Button) {
		account := accountInput.Text()
		wxSecret := secretInput.Text()
		if len(account) == 0 {
			ui.MsgBoxError(w, "错误提示", "手机号不能为空")
			return
		}
		if len(wxSecret) == 0 {
			ui.MsgBoxError(w, "错误提示", "内容不能为空")
			return
		}
		ok := utils.SendMessage(account, `{"code":"`+wxSecret+`"}`, "SMS_200185703", "冰点")
		if ok {
			ui.MsgBox(w, "提示", "发送成功")
		} else {
			ui.MsgBox(w, "提示", "发送异常")
		}
	})
	return hbox
}
