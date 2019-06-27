package main

/*
Пример запроса на получение токена:
https://oauth.vk.com/authorize?client_id=2890984&scope=market,notify,photos,friends,audio,video,notes,pages,docs,status,questions,offers,wall,groups,messages,notifications,stats,ads,offline&redirect_uri=http://api.vk.com/blank.html&display=page&response_type=token
*/

import (
	"github.com/salaleser/vk-poster/command"
	"github.com/salaleser/vk-poster/gui"

	"github.com/google/gxui/drivers/gl"
)

func main() {
	command.LoadCommands()
	gl.StartDriver(gui.MainWindow)
}
