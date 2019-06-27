package command

import (
	"log"
	"os"
	"path"
	"strconv"
	"strings"

	"github.com/salaleser/vk-api/method/market"
	"github.com/salaleser/vk-api/method/photos"
	"github.com/salaleser/vk-api/util"
	"github.com/salaleser/vk-poster/logger"
)

// AddProduct загружает данные из указанного файла в сообщество
func addProduct(groupID string, categoryID string, albumID string,
	imageURL string, title string, data string, price string) {
	url := strings.Split(imageURL, "/")
	if len(url) < 3 {
		logger.Appendf("Ошибка при попытке разобрать адрес изображения %q", imageURL)
		return
	}
	filepath := path.Join("temp", url[len(url)-2])
	err := os.MkdirAll(filepath, 0755)
	if err != nil {
		log.Fatal(err)
	}
	fullname := path.Join(filepath, url[len(url)-1])

	logger.Appendf("Загружаю файл с адреса %q", imageURL)
	util.DownloadFile(imageURL, fullname)
	logger.Appendf("Файл загружен в %q", fullname)

	uo := photos.GetMarketUploadServer(groupID, "1", "", "", "")

	fo := util.UploadFile(uo.R.UploadURL, fullname)
	photo := fo.Photo
	server := strconv.Itoa(fo.Server)
	hash := fo.Hash
	cropData := fo.CropData
	cropHash := fo.CropHash

	mf := photos.SaveMarketPhoto(groupID, photo, server, hash, cropData, cropHash)
	if len(mf.R) == 0 {
		logger.Appendf("Не удалось сохранить файл %q на сервере", fullname)
		return
	}
	ownerID := strconv.Itoa(mf.R[0].OwnerID)
	name := title
	description := data
	cat := categoryID
	p := parsePrice(price)
	deleted := "0"
	mainPhotoID := strconv.Itoa(mf.R[0].ID)
	photoIDs := ""

	ao := market.Add(ownerID, name, description, cat, p, deleted, mainPhotoID, photoIDs)
	itemID := ao.R.MarketItemID
	if itemID == 0 {
		logger.Appendf("Ошибка при попытке добавить товар в сообщество!")
		return
	}
	logger.Appendf("Товар добавлен успешно, вот его идентификатор: %d", itemID)

	atao := market.AddToAlbum(groupID, strconv.Itoa(itemID), albumID)
	if atao.Response == 0 {
		logger.Appendf("Произошла ошибка при попытке добавить товар в подборку!")
		return
	}
	logger.Appendf("Товар добавлен в альбом успешно")
}

func parsePrice(price string) string {
	var f string
	for i := 0; i < len(price); i++ {
		l := price[i : i+1]
		_, err := strconv.Atoi(l)
		if err != nil {
			break
		}
		f += l
	}
	return f
}
