package main

import (
	"fmt"

	"github.com/andlabs/ui"
	_ "github.com/andlabs/ui/winmanifest"
	"github.com/astaxie/beego/toolbox"
	"github.com/sirupsen/logrus"

	"github.com/libra412/auto_recharge/services"
)

var log = logrus.New()

//
func main() {
	tk := toolbox.NewTask("tk", "0 * * * * ?", f)
	toolbox.AddTask("tk", tk)
	ui.Main(setUp)
}

//
func setUp() {
	// 初始化窗口
	mainwin := ui.NewWindow("挂机软件V0.0.2.beta", 640, 480, false)
	mainwin.OnClosing(func(*ui.Window) bool {
		ui.Quit()
		return true
	})
	ui.OnShouldQuit(func() bool {
		mainwin.Destroy()
		return true
	})

	// 新建tab
	tab := ui.NewTab()
	tab.Append("火币", makeControl(mainwin))
	tab.SetMargined(0, true)
	// 设置tab页
	mainwin.SetChild(tab)
	mainwin.SetMargined(true)

	// 最后显示
	mainwin.Show()
}

// 界面布局
func makeControl(w *ui.Window) ui.Control {
	hbox := ui.NewHorizontalBox()
	hbox.SetPadded(true)

	entryForm := ui.NewForm()
	entryForm.SetPadded(true)
	hbox.Append(entryForm, false)
	// 设置发送请求按钮
	isRequest := false
	requestButton := ui.NewButton("开始接单")
	entryForm.Append("", requestButton, false)
	valueInput := ui.NewEntry()
	entryForm.Append("均价", valueInput, false)
	timeInput := ui.NewEntry()
	entryForm.Append("频率", timeInput, false)

	// 开始操作按钮
	requestButton.OnClicked(func(*ui.Button) {

		if isRequest {
			isRequest = false
			toolbox.StopTask()
			requestButton.SetText("开始接单")
		} else {
			isRequest = true
			value := valueInput.Text()
			timeValue := timeInput.Text()
			if len(value) == 0 {
				ui.MsgBoxError(w, "错误提示", "点位不能为空")
				return
			}
			if len(timeValue) == 0 {
				ui.MsgBoxError(w, "错误提示", "频率不能为空")
				return
			}
			toolbox.StartTask()
			requestButton.SetText("停止接单")
		}
	})
	return hbox
}

// 自动接单
func f() error {

	data := services.GetKLine("BTC_CQ", "4hour", 100)
	fmt.Println("请求数据", data)
	count := len(data.Data)
	for i := 0; i < count; i++ {
		fmt.Println(data.Data[i])
	}
	return nil
}
