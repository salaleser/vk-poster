package logger

import (
	"fmt"
	"log"

	"github.com/google/gxui"
)

// Field объект с логом
var Field gxui.CodeEditor

// Append добавляет текст в лог
func Append(s string) {
	log.Println(s)
	if Field.Text() != "" {
		s = "\n" + s
	}
	Field.SetText(Field.Text() + s)
}

// Appendf добавляет форматированный текст в лог
func Appendf(f string, a ...interface{}) {
	s := fmt.Sprintf(f, a...)
	log.Println(s)
	if Field.Text() != "" {
		s = "\n" + s
	}
	Field.SetText(Field.Text() + s)
}

// Clear очищает лог
func Clear() {
	Field.SetText("")
}
