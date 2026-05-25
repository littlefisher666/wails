package main

import (
	"embed"
	goruntime "runtime"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/menu"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
)

//go:embed all:frontend/dist
var assets embed.FS

func main() {
	// 只创建一个 App 实例，便于前端多次调用时共享同一份后端状态。
	app := NewApp()

	err := wails.Run(&options.App{
		Title:  "helloworld",
		Width:  1024,
		Height: 768,
		AssetServer: &assetserver.Options{
			Assets: assets,
		},
		BackgroundColour:  &options.RGBA{R: 27, G: 38, B: 54, A: 1},
		HideWindowOnClose: true,
		Menu:              newApplicationMenu(app),
		OnStartup:         app.startup,
		Bind: []interface{}{
			// 绑定后，App 的公开方法会生成到 frontend/wailsjs 供 Vue 调用。
			app,
		},
	})

	if err != nil {
		println("Error:", err.Error())
	}
}

func newApplicationMenu(app *App) *menu.Menu {
	applicationMenu := defaultApplicationMenu()
	windowMenu := applicationMenu.AddSubmenu("窗口控制")

	// 菜单回调不返回错误，窗口命令内部会先校验 Wails 上下文是否就绪。
	windowMenu.AddText("显示窗口", nil, func(_ *menu.CallbackData) {
		_, _ = app.ShowWindow()
	})
	windowMenu.AddText("隐藏到后台", nil, func(_ *menu.CallbackData) {
		_, _ = app.HideToTray()
	})
	windowMenu.AddText("窗口居中", nil, func(_ *menu.CallbackData) {
		_, _ = app.CenterWindow()
	})
	windowMenu.AddSeparator()
	windowMenu.AddText("退出应用", nil, func(_ *menu.CallbackData) {
		_ = app.QuitApp()
	})

	return applicationMenu
}

func defaultApplicationMenu() *menu.Menu {
	switch goruntime.GOOS {
	case "darwin":
		return menu.NewMenuFromItems(menu.AppMenu(), menu.EditMenu(), menu.WindowMenu())
	case "windows":
		return menu.NewMenuFromItems(menu.EditMenu(), menu.WindowMenu())
	default:
		return menu.NewMenu()
	}
}
