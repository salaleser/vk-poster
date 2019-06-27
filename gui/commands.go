package gui

import (
	"path/filepath"

	"github.com/salaleser/poster/command"
	"github.com/salaleser/poster/logger"
	"github.com/salaleser/poster/parser"
	"github.com/salaleser/vk-api/util"

	"github.com/google/gxui"
)

var (
	title        gxui.Label          // Наименование команды
	description  gxui.Label          // Описание команды
	filename     gxui.TextBox        // Поле с именем файла
	titles       []gxui.Label        // Наименования аргументов
	descriptions []gxui.Label        // Описания аргументов
	boxes        []gxui.TextBox      // Текстовые поля аргументов
	layouts      []gxui.LinearLayout // Контейнеры для аргументов

	dropList gxui.DropDownList // Перечень команд
)

var (
	titleColor       = gxui.Color{R: 1.0, G: 0.8, B: 0.3, A: 1}
	descriptionColor = gxui.Color{R: 0.3, G: 1.0, B: 0.5, A: 1}
	paramTitleColor  = gxui.Color{R: 0.7, G: 0.6, B: 0.9, A: 1}
	tipColor         = gxui.Color{R: 0.6, G: 0.6, B: 0.6, A: 1}
)

func methodPicker(theme gxui.Theme, overlay gxui.BubbleOverlay) gxui.Control {
	adapter := gxui.CreateDefaultAdapter()
	adapter.SetItems(command.Items)

	layout := theme.CreateLinearLayout()
	layout.SetDirection(gxui.TopToBottom)

	methodLayout := theme.CreateLinearLayout()
	methodLayout.SetDirection(gxui.LeftToRight)

	label := theme.CreateLabel()
	label.SetText("Команда:")
	methodLayout.AddChild(label)

	dropList = theme.CreateDropDownList()
	dropList.SetAdapter(adapter)
	dropList.SetBubbleOverlay(overlay)
	dropList.Select(command.Items[0])
	methodLayout.AddChild(dropList)

	tipLabel := theme.CreateLabel()
	tipLabel.SetColor(tipColor)
	tipLabel.SetText("Если будет указан файл, то программа попытается брать аргументы из него, " +
		"иначе данные будут извлечены из текстовых полей")
	layout.AddChild(tipLabel)
	layout.AddChild(methodLayout)

	title = theme.CreateLabel()
	title.SetColor(titleColor)
	layout.AddChild(title)

	description = theme.CreateLabel()
	description.SetColor(descriptionColor)
	layout.AddChild(description)

	filenameLayout := theme.CreateLinearLayout()
	filenameLayout.SetDirection(gxui.LeftToRight)

	filenameLabel := theme.CreateLabel()
	filenameLabel.SetText("Имя файла:")
	filename = theme.CreateTextBox()
	filename.SetDesiredWidth(width - 100 - filenameLabel.Size().W)
	filenameLayout.AddChild(filenameLabel)
	filenameLayout.AddChild(filename)
	layout.AddChild(filenameLayout)

	for i := 0; i < 10; i++ {
		titles = append(titles, theme.CreateLabel())
		boxes = append(boxes, theme.CreateTextBox())
		descriptions = append(descriptions, theme.CreateLabel())
		layouts = append(layouts, theme.CreateLinearLayout())
		layouts[i].SetDirection(gxui.LeftToRight)
		layouts[i].SetVisible(false)
		layouts[i].AddChild(boxes[i])
		layouts[i].AddChild(titles[i])
		titles[i].SetColor(paramTitleColor)
		layouts[i].AddChild(descriptions[i])
		layout.AddChild(layouts[i])
	}

	logButtonsLayout := theme.CreateLinearLayout()
	logButtonsLayout.SetDirection(gxui.LeftToRight)

	run := theme.CreateButton()
	run.SetText("Выполнить")
	run.OnClick(func(gxui.MouseEvent) {
		if tokenBox.Text() == "" {
			logger.Append("Для выполнения запроса указать токен вручную или получить по кнопке")
			return
		}
		util.UserToken = tokenBox.Text()
		runCommand(dropList.Selected().(string))
	})
	logButtonsLayout.AddChild(run)

	clear := theme.CreateButton()
	clear.SetText("Очистить")
	clear.OnClick(func(gxui.MouseEvent) {
		logger.Clear()
	})
	logButtonsLayout.AddChild(clear)

	layout.AddChild(logButtonsLayout)

	logger.Field = theme.CreateCodeEditor()
	logger.Field.SetDesiredWidth(width - 40)
	logLayout := theme.CreateScrollLayout()
	logLayout.SetChild(logger.Field)
	layout.AddChild(logLayout)

	update(command.GetCommand(command.Items[0]))
	dropList.OnSelectionChanged(func(item gxui.AdapterItem) {
		update(command.GetCommand(dropList.Selected().(string)))
	})

	return layout
}

func runCommand(cmd string) {
	c := command.GetCommand(cmd)

	a := []parser.Args{
		parser.Args{
			Arg1:  boxes[0].Text(),
			Arg2:  boxes[1].Text(),
			Arg3:  boxes[2].Text(),
			Arg4:  boxes[3].Text(),
			Arg5:  boxes[4].Text(),
			Arg6:  boxes[5].Text(),
			Arg7:  boxes[6].Text(),
			Arg8:  boxes[7].Text(),
			Arg9:  boxes[8].Text(),
			Arg10: boxes[9].Text(),
		},
	}
	f := filename.Text()
	if f != "" {
		if filepath.Ext(f) != ".csv" {
			logger.Append("Ошибка: файл должен быть CSV-формата")
			return
		}
		a = parser.ParseCSV(f, a)
	}
	c.Func(a)
}
