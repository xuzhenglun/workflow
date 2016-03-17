package cli

import (
	"fmt"
	"os"
	"strings"

	"github.com/xuzhenglun/workflow/core"
)

var coreIoBus core.CoreIoBus

func Run(c core.CoreIoBus) {
	cmd := initCmd()
	fmt.Scanln()
	fmt.Println("Command Line Enabled.")
	coreIoBus = c
	for {
		fmt.Printf(">>")
		c := ScanLine()
		input := strings.Split(c, " ")
		if input[0] != "" {
			if f, ok := (*cmd)[input[0]]; ok {
				f(input[1:]...)
			} else {
				fmt.Println("Error: No such command \"" + input[0] + "\"")
			}
		}
	}
}

func initCmd() *map[string]func(param ...string) string {
	cmdTable := make(map[string]func(param ...string) string)
	cmdTable["exit"] = exit
	cmdTable["echo"] = echo
	cmdTable["lua"] = runlua
	return &cmdTable
}

func exit(param ...string) string {
	os.Exit(0)
	fmt.Println("bye")
	return "SUCC"
}

func echo(param ...string) string {
	for _, v := range param {
		if v != "" {
			fmt.Print(v + " ")
		}
	}
	fmt.Println()
	return "SUCC"
}

func runlua(param ...string) string {
	for _, v := range param {
		coreIoBus.RawLuaHandler(v)
	}
	return "Done"
}

func ScanLine() string {
	var c byte
	var err error
	var b []byte
	for err == nil {
		_, err = fmt.Scanf("%c", &c)
		if c != '\n' {
			b = append(b, c)
		} else {
			break
		}
	}
	return string(b)
}
