package main

import (
	"fmt"
	"os"
	"time"
	"vtools/capture"
	"vtools/tools"

	"github.com/therecipe/qt/core"
	"github.com/therecipe/qt/gui"
	"github.com/therecipe/qt/widgets"
	"golang.design/x/hotkey"
	"golang.design/x/hotkey/mainthread"
)

var (
	recognizeFlag = false
)

type MySystrayMenu struct {
	widgets.QSystemTrayIcon
	_ func()       `signal:"recognize"`
	_ func(string) `signal:"setIconFromSignal"`
}

func main() {
	go tools.RunOcr()
	core.QCoreApplication_SetAttribute(core.Qt__AA_EnableHighDpiScaling, true)
	app := widgets.NewQApplication(len(os.Args), os.Args)
	app.SetQuitOnLastWindowClosed(false)
	clipboard := app.Clipboard()
	window := widgets.NewQMainWindow(nil, 0)
	recognizeFlagPtr := &recognizeFlag
	systray := NewMySystrayMenu(nil)
	onOcrData := func(text string) {
		systray.SetIconFromSignal("idle.ico")
		systray.SetToolTip("服务空闲中...")
		if text != "" {
			clipboard.SetText(text, gui.QClipboard__Clipboard)
			systray.ShowMessage("截图识别", fmt.Sprintf("识别成功,请粘贴:\n%s", text), widgets.QSystemTrayIcon__Information, 5000)
		}
		*recognizeFlagPtr = false
	}
	bulinbulin := func() {
		// for {
		if *recognizeFlagPtr {
			systray.SetIconFromSignal("busy.ico")
			time.Sleep(time.Millisecond * 500)
			// systray.SetIconFromSignal("idle.ico")
			// fmt.Println("bulinbulin")
		} else {
			return
		}
		// }
	}
	recognize := func() {
		if !*recognizeFlagPtr {
			systray.SetToolTip("识别中...")
			*recognizeFlagPtr = true
			go bulinbulin()
			sniper := capture.NNewSniper(window, 0, onOcrData)
			sniper.Recognize()
		} else {
			// systray.ShowMessage("截图识别", fmt.Sprintf("识别中请勿重复操作"), widgets.QSystemTrayIcon__Warning, 5000)
		}
	}
	systray.ConnectSetIconFromSignal(func(v0 string) {
		fmt.Print("setIconFromSignal:", v0)
		systray.SetIcon(gui.NewQIcon5(v0))
	})
	systray.ConnectRecognize(recognize)
	systray.SetIcon(gui.NewQIcon5("idle.ico"))
	systray.SetToolTip("服务空闲中...")
	systrayMenu := widgets.NewQMenu(nil)

	systrayMenu.AddAction("截图识别").ConnectTriggered(func(checked bool) {
		recognize()
	}) 
	systrayMenu.AddAction("退出").ConnectTriggered(func(checked bool) {
		fmt.Println("退出")
		tools.KillOcr()
		app.Quit()
	})
	systray.SetContextMenu(systrayMenu)
	systray.Show()
	go func() {
		mainthread.Init(func() {
			hk := hotkey.New([]hotkey.Modifier{hotkey.ModCtrl, hotkey.ModShift}, hotkey.KeyS)
			err := hk.Register()
			if err != nil {
				return
			}
			fmt.Printf("hotkey: %v is registered\n", hk)

			for {
				select {
				case <-hk.Keydown():
					fmt.Println("hotkey:", hk, "is pressed")
					go systray.Recognize()
				}
			}
		})
	}()

	os.Exit(app.Exec())
}
