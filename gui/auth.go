package gui

import (
	"log"

	"github.com/google/gxui"
	"github.com/salaleser/vk-api/method"
	"github.com/salaleser/vk-poster/command"
	"github.com/salaleser/vk-poster/logger"
)

var (
	clientBox   gxui.TextBox
	loginBox    gxui.TextBox
	passwordBox gxui.TextBox
	tokenBox    gxui.TextBox
)

func authPicker(theme gxui.Theme) gxui.Control {
	adapter := gxui.CreateDefaultAdapter()
	adapter.SetItems(command.Items)

	layout := theme.CreateLinearLayout()
	layout.SetDirection(gxui.TopToBottom)

	clientLayout := theme.CreateLinearLayout()
	clientLayout.SetDirection(gxui.LeftToRight)

	clientLabel := theme.CreateLabel()
	clientLabel.SetText("ID приложения:")
	clientLayout.AddChild(clientLabel)

	clientBox = theme.CreateTextBox()
	clientBox.SetDesiredWidth(60)
	clientLayout.AddChild(clientBox)

	salaleserClientID := theme.CreateButton()
	salaleserClientID.SetText("salaleser")
	salaleserClientID.OnClick(func(gxui.MouseEvent) {
		clientBox.SetText("6713750")
	})
	clientLayout.AddChild(salaleserClientID)

	iPhoneClientID := theme.CreateButton()
	iPhoneClientID.SetText("iPhone")
	iPhoneClientID.OnClick(func(gxui.MouseEvent) {
		clientBox.SetText("3140623")
	})
	clientLayout.AddChild(iPhoneClientID)

	iPadClientID := theme.CreateButton()
	iPadClientID.SetText("iPad")
	iPadClientID.OnClick(func(gxui.MouseEvent) {
		clientBox.SetText("3682744")
	})
	clientLayout.AddChild(iPadClientID)

	androidClientID := theme.CreateButton()
	androidClientID.SetText("Android")
	androidClientID.OnClick(func(gxui.MouseEvent) {
		clientBox.SetText("2274003")
	})
	clientLayout.AddChild(androidClientID)

	windowsClientID := theme.CreateButton()
	windowsClientID.SetText("Windows Phone")
	windowsClientID.OnClick(func(gxui.MouseEvent) {
		clientBox.SetText("3502557")
	})
	clientLayout.AddChild(windowsClientID)

	layout.AddChild(clientLayout)

	loginLayout := theme.CreateLinearLayout()
	loginLayout.SetDirection(gxui.LeftToRight)

	loginLabel := theme.CreateLabel()
	loginLabel.SetText("Логин:")
	loginLayout.AddChild(loginLabel)

	loginBox = theme.CreateTextBox()
	loginBox.SetDesiredWidth(120)
	loginLayout.AddChild(loginBox)

	layout.AddChild(loginLayout)

	passwordLayout := theme.CreateLinearLayout()
	passwordLayout.SetDirection(gxui.LeftToRight)

	passwordLabel := theme.CreateLabel()
	passwordLabel.SetText("Пароль:")
	passwordLayout.AddChild(passwordLabel)

	passwordBox = theme.CreateTextBox()
	passwordBox.SetDesiredWidth(100)
	passwordLayout.AddChild(passwordBox)

	layout.AddChild(passwordLayout)

	buttonsLayout := theme.CreateLinearLayout()
	buttonsLayout.SetDirection(gxui.LeftToRight)

	login := theme.CreateButton()
	login.SetText("Получить токен")
	login.OnClick(func(gxui.MouseEvent) {
		c := clientBox.Text()
		if c == "" {
			logger.Append("Для получения токена необходимо указать ID приложения.\n" +
				"Для этого можно использовать один из предложенных\n" +
				"вариантов или создать свой по этой инструкции\n" +
				"https://readd.org/kak-poluchit-access_token-vkontakte/")
			return
		}

		l := loginBox.Text()
		p := passwordBox.Text()
		if l == "" || p == "" {
			logger.Append("Для получения токена необходимо указать логин и пароль")
			return
		}

		token, err := method.Auth(l, p, c)
		if err != nil {
			log.Println(err)
			return
		}

		tokenBox.SetText(token)
	})
	buttonsLayout.AddChild(login)
	layout.AddChild(buttonsLayout)

	tokenLayout := theme.CreateLinearLayout()
	tokenLayout.SetDirection(gxui.LeftToRight)

	tokenLabel := theme.CreateLabel()
	tokenLabel.SetText("Токен:")
	tokenLayout.AddChild(tokenLabel)

	tokenBox = theme.CreateTextBox()
	tokenBox.SetDesiredWidth(width - tokenLabel.Size().W)
	tokenLayout.AddChild(tokenBox)

	layout.AddChild(tokenLayout)

	return layout
}
