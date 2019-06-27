package command

import (
	"strconv"

	"github.com/salaleser/vk-api/method/groups"
	"github.com/salaleser/vk-api/method/market"
	"github.com/salaleser/vk-api/method/wall"
	"github.com/salaleser/vk-poster/logger"
	"github.com/salaleser/vk-poster/parser"
)

const defaultOwnerID = "171524015"

// Command описывает команду
type Command struct {
	Func        func([]parser.Args)
	Title       string
	Description string
	Return      string
	Params      []Param
}

// Param описывает один из аргументов в команде
type Param struct {
	Title       string
	Description string
	Defaults    string
}

// Commands содержит список команд
var Commands = []Command{}

// Items содержит список наименований команд
var Items = []string{}

// LoadCommands загружает команды
func LoadCommands() {
	Commands = append(Commands,
		Command{
			Title: "Добавить товар",
			Func: func(args []parser.Args) {
				o := market.AddAlbum(args[0].GroupID, args[0].AlbumTitle, "", "0")
				logger.Appendf("Подборка %q добавлена", args[0].AlbumTitle)
				for _, a := range args {
					addProduct(a.GroupID, a.CategoryID, strconv.Itoa(o.R.MarketAlbumID),
						a.Image, a.Title, "Артикул: "+a.Model+", "+a.Data, a.Price)
				}
			},
			Description: "Для добавления товаров в сообщество необходимо указать файл в " +
				"соответствующее поле на этой вкладке. Эта команда создаст альбом и " +
				"наполнит его товарами.",
			Params: []Param{
				Param{
					Title:       "group_id",
					Description: "идентификатор или короткое имя сообщества.",
					Defaults:    defaultOwnerID,
				},
				Param{
					Title:       "album_title",
					Description: "наименование подборки",
					Defaults:    "<имя_подборки>",
				},
				Param{
					Title:       "category_id",
					Description: "идентификационный номер категории",
					Defaults:    "5",
				},
			},
			Return: "После успешного выполнения возвращает объект, содержащий число результатов в поле count (integer) и массив идентификаторов пользователей в поле items ([integer]).",
		})
	Commands = append(Commands,
		Command{
			Title: "Удалить товары",
			Func: func(args []parser.Args) {
				for _, a := range args {
					limit, err := strconv.Atoi(a.Arg2)
					if err != nil {
						logger.Appendf("Ошибка: %s", err)
						continue
					}
					deleteProducts(a.Arg1, limit)
				}
			},
			Description: "Удаляет товары в сообществе.",
			Params: []Param{
				Param{
					Title:       "owner_id",
					Description: "идентификатор или короткое имя сообщества.",
					Defaults:    defaultOwnerID,
				},
				Param{
					Title:       "limit",
					Description: "ограничение на количество удаляемых товаров за одну сессию",
					Defaults:    "50",
				},
			},
		})
	Commands = append(Commands,
		Command{
			Title: "groups.getMembers",
			Func: func(args []parser.Args) {
				for _, a := range args {
					o := groups.GetMembers(a.Arg1, a.Arg2, a.Arg3, a.Arg4, a.Arg5, a.Arg6)
					logger.Appendf("Всего в сообществе %d участников", o.R.Count)
					for _, item := range o.R.Items {
						logger.Appendf("Идентификатор участника: %d", item)
					}
				}
			},
			Description: "Возвращает список участников сообщества.",
			Params: []Param{
				Param{
					Title:       "group_id",
					Description: "идентификатор или короткое имя сообщества.",
					Defaults:    defaultOwnerID,
				},
				Param{
					Title:       "offset",
					Description: "смещение, необходимое для выборки определенного подмножества участников. По умолчанию 0.",
				},
				Param{
					Title:       "sort",
					Description: "сортировка, с которой необходимо вернуть список участников.",
				},
				Param{
					Title:       "count",
					Description: "количество участников сообщества, информацию о которых необходимо получить.",
				},
				Param{
					Title:       "fields",
					Description: "список дополнительных полей, которые необходимо вернуть.",
				},
				Param{
					Title:       "filter",
					Description: "(фильтр)",
				},
			},
			Return: "После успешного выполнения возвращает объект, содержащий число результатов в поле count (integer) и массив идентификаторов пользователей в поле items ([integer]).",
		})
	Commands = append(Commands,
		Command{
			Title: "wall.post",
			Func: func(args []parser.Args) {
				for _, a := range args {
					o := wall.Post(a.Arg1, a.Arg2, a.Arg3, a.Arg4)
					logger.Appendf("Идентификатор записи: %d", o.R.PostID)
				}
			},
			Description: "Позволяет создать запись на стене, предложить запись на стене публичной страницы, опубликовать существующую отложенную запись.",
			Params: []Param{
				Param{
					Title:       "owner_id",
					Description: "идентификатор пользователя или сообщества, на стене которого должна быть опубликована запись.",
					Defaults:    defaultOwnerID,
				},
				Param{
					Title:       "friends_only",
					Description: "1 — запись будет доступна только друзьям, 0 — всем пользователям. По умолчанию публикуемые записи доступны всем пользователям.",
					Defaults:    "0",
				},
				Param{
					Title:       "from_group",
					Description: "данный параметр учитывается, если owner_id < 0 (запись публикуется на стене группы). 1 — запись будет опубликована от имени группы, 0 — запись будет опубликована от имени пользователя (по умолчанию).",
				},
				Param{
					Title:       "message",
					Description: "текст сообщения (является обязательным, если не задан параметр attachments)",
				},
			},
			Return: "После успешного выполнения возвращает идентификатор созданной записи (post_id).",
		})
	Commands = append(Commands,
		Command{
			Title: "market.get",
			Func: func(args []parser.Args) {
				for _, a := range args {
					o := market.Get(a.Arg1, a.Arg2, a.Arg3, a.Arg4, a.Arg5)
					logger.Appendf("Количество товаров в сообществе: %d", o.R.Count)
					for _, item := range o.R.Items {
						logger.Appendf("%s (%d)", item.Title, item.ID)
					}
				}
			},
			Description: "Возвращает список товаров в сообществе.",
			Params: []Param{
				Param{
					Title:       "owner_id",
					Description: "идентификатор владельца товаров.",
					Defaults:    defaultOwnerID,
				},
				Param{
					Title:       "album_id",
					Description: "идентификатор подборки, товары из которой нужно вернуть.",
				},
				Param{
					Title:       "count",
					Description: "количество возвращаемых товаров.",
				},
				Param{
					Title:       "offset",
					Description: "смещение относительно первого найденного товара для выборки определенного подмножества.",
				},
				Param{
					Title:       "extended",
					Description: "1 — будут возвращены дополнительные поля likes, can_comment, can_repost, photos, views_count. По умолчанию эти поля не возвращается.",
				},
			},
			Return: "После успешного выполнения возвращает объект, содержащий число результатов в поле count и массив объектов товаров в поле items.",
		})
	Commands = append(Commands,
		Command{
			Title: "market.getAlbums",
			Func: func(args []parser.Args) {
				for _, a := range args {
					o := market.GetAlbums(a.Arg1, a.Arg2, a.Arg3)
					logger.Appendf("Количество подборок с товарами: %d", o.R.Count)
					for _, item := range o.R.Items {
						logger.Appendf("%s (%d)", item.Title, item.ID)
					}
				}
			},
			Description: "Возвращает список подборок с товарами.",
			Params: []Param{
				Param{
					Title:       "owner_id",
					Description: "идентификатор владельца товаров.",
					Defaults:    defaultOwnerID,
				},
				Param{
					Title:       "offset",
					Description: "смещение относительно первой найденной подборки для выборки определенного подмножества.",
				},
				Param{
					Title:       "count",
					Description: "количество возвращаемых товаров.",
				},
			},
			Return: "После успешного выполнения возвращает объект, содержащий число результатов в поле count и массив объектов подборок в поле items.",
		})
	Commands = append(Commands,
		Command{
			Title: "market.delete",
			Func: func(args []parser.Args) {
				for _, a := range args {
					o := market.Delete(a.Arg1, a.Arg2)
					if o.R == 1 {
						logger.Append("Товар успешно удален")
					} else {
						logger.Append("Товар не удален!")
					}
				}
			},
			Description: "Удаляет товар.",
			Params: []Param{
				Param{
					Title:       "owner_id",
					Description: "идентификатор владельца товаров.",
					Defaults:    defaultOwnerID,
				},
				Param{
					Title:       "item_id",
					Description: "идентификатор товара.",
				},
			},
			Return: "После успешного выполнения возвращает 1.",
		})
	Commands = append(Commands,
		Command{
			Title: "market.deleteAlbum",
			Func: func(args []parser.Args) {
				for _, a := range args {
					o := market.DeleteAlbum(a.Arg1, a.Arg2)
					if o.R == 1 {
						logger.Append("Подборка успешно удалена")
					} else {
						logger.Append("Подборка не удалена!")
					}
				}
			},
			Description: "Удаляет подборку с товарами.",
			Params: []Param{
				Param{
					Title:       "owner_id",
					Description: "идентификатор владельца подборки.",
					Defaults:    defaultOwnerID,
				},
				Param{
					Title:       "album_id",
					Description: "идентификатор подборки.",
				},
			},
			Return: "После успешного выполнения возвращает 1.",
		})
	Commands = append(Commands,
		Command{
			Title: "market.getCategories",
			Func: func(args []parser.Args) {
				for _, a := range args {
					o := market.GetCategories(a.Arg1, a.Arg2)
					logger.Appendf("Всего категорий: %d", o.R.Count)
					for _, item := range o.R.Items {
						logger.Appendf("%s (%d)", item.Name, item.ID)
					}
				}
			},
			Description: "Возвращает список категорий для товаров.",
			Params: []Param{
				Param{
					Title:       "count",
					Description: "количество категорий, информацию о которых необходимо вернуть.",
					Defaults:    "20",
				},
				Param{
					Title:       "offset",
					Description: "смещение, необходимое для выборки определенного подмножества категорий.",
				},
			},
			Return: "После успешного выполнения возвращает список объектов category.",
		})

	for _, c := range Commands {
		Items = append(Items, c.Title)
	}
}

// GetCommand возвращает команду по наименованию
func GetCommand(title string) Command {
	for _, c := range Commands {
		if c.Title == title {
			return c
		}
	}
	return Command{}
}
