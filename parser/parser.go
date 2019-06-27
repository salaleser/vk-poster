package parser

import (
	"bufio"
	"encoding/csv"
	"io"
	"os"

	"github.com/salaleser/vk-poster/logger"
)

// Args описывает объект с аргументами
type Args struct {
	AlbumTitle string
	GroupID    string
	CategoryID string
	Title      string
	Image      string
	Price      string
	Model      string
	Data       string

	Arg1  string
	Arg2  string
	Arg3  string
	Arg4  string
	Arg5  string
	Arg6  string
	Arg7  string
	Arg8  string
	Arg9  string
	Arg10 string
}

// ParseCSV создает из CSV-файла объекты
// FIXME a — постоянный временный костыль
func ParseCSV(filename string, a []Args) []Args {
	file, err := os.Open(filename)
	if err != nil {
		logger.Appendf("Ошибка при попытке прочитать файл %q (%s)", filename, err)
		return []Args{}
	}
	reader := csv.NewReader(bufio.NewReader(file))
	var args []Args
	for {
		line, err := reader.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			logger.Appendf("Ошибка при попытке разобрать файл %q (%s)", filename, err)
			return []Args{}
		}
		args = append(args,
			Args{
				Title: line[0],
				Image: line[1],
				Price: line[2],
				Model: line[3],
				Data:  line[4],

				GroupID:    a[0].Arg1,
				AlbumTitle: a[0].Arg2,
				CategoryID: a[0].Arg3,
			})
	}
	return args
}
