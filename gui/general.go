package gui

import (
	"github.com/google/gxui"
	"github.com/google/gxui/math"
	"github.com/google/gxui/samples/flags"
	"github.com/salaleser/vk-poster/command"
	"github.com/salaleser/vk-poster/logger"
)

const (
	width  = 455
	height = 800
)

var filenameText string

// MainWindow создает основное окно
func MainWindow(driver gxui.Driver) {
	theme := flags.CreateTheme(driver)

	window := theme.CreateWindow(width, height, "Poster")
	window.SetPosition(math.Point{X: 1366 - width, Y: 0})

	overlay := theme.CreateBubbleOverlay()

	holder := theme.CreatePanelHolder()
	holder.AddPanel(methodPicker(theme, overlay), "Выбор команды")
	holder.AddPanel(filePicker(theme), "Выбор CSV-файла для разбора аргументов")
	holder.AddPanel(authPicker(theme), "Авторизация")

	holder.Tab(0).OnClick(func(gxui.MouseEvent) {
		filename.SetText(filenameText)
	})

	window.OnKeyDown(func(ev gxui.KeyboardEvent) {
		if ev.Key == gxui.KeyEnter || ev.Key == gxui.KeyKpEnter {
			runCommand(dropList.Selected().(string))
		}
	})

	window.OnKeyDown(func(ev gxui.KeyboardEvent) {
		if ev.Key == gxui.KeyBackspace || ev.Key == gxui.KeyDelete {
			logger.Clear()
		}
	})

	window.SetScale(flags.DefaultScaleFactor)
	window.AddChild(holder)
	window.AddChild(overlay)
	window.OnClose(driver.Terminate)
	window.SetPadding(math.Spacing{L: 10, T: 10, R: 10, B: 10})
}

func update(c command.Command) {
	title.SetText(c.Title)
	description.SetText(c.Description)
	for i := 0; i < len(c.Params); i++ {
		layouts[i].SetVisible(true)
		titles[i].SetText(c.Params[i].Title)
		descriptions[i].SetText(c.Params[i].Description)
		boxes[i].SetText(c.Params[i].Defaults)
	}
	for i := len(c.Params); i < len(boxes); i++ {
		layouts[i].SetVisible(false)
	}
}
