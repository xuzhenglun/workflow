package api

import (
	"log"

	"github.com/yuin/gopher-lua"
)

var List = map[string]func(*lua.LState) int{
	"log": __log__,
}

func __log__(l *lua.LState) int {
	v := l.ToString(1)
	log.Println(v)
	return 0
}
