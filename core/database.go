package core

import (
	"encoding/json"
	"log"
	"strconv"

	"github.com/yuin/gopher-lua"
)

func (this VMs) InitDatabase(l *lua.LState, a *Activity) {
	l.SetGlobal("AddRow", l.NewFunction(this.AddRowFunc(*a)))
	l.SetGlobal("DelRow", l.NewFunction(this.DelRowFunc()))
	l.SetGlobal("ModRow", l.NewFunction(this.ModRowFunc(*a)))
	l.SetGlobal("FindRow", l.NewFunction(this.FindRowFunc()))
}

func (this VMs) AddRowFunc(a Activity) func(l *lua.LState) int {
	notnull := this.Re.FindAllString(a.Helper, -1)
	return func(l *lua.LState) int {
		if str, ok := l.Get(-1).(lua.LString); ok {
			l.Pop(1)
			log.Println("--->" + string(str))
			args := make(map[string]interface{})
			err := json.Unmarshal([]byte(str), &args)
			if err != nil {
				log.Println("Error: Json Decode fail\n" + string(str))
			}
			if err := this.Db.AddRow(notnull, args); err != nil {
				log.Println(err)
				l.Push(lua.LString("FAIL"))
				return 1
			} else {
				l.Push(lua.LString("SUCC"))
				return 1
			}
		} else {
			l.Push(lua.LString("ERR"))
			return 1
		}
	}
}

func (this VMs) ModRowFunc(a Activity) func(l *lua.LState) int {
	notnull := this.Re.FindAllString(a.Helper, -1)
	return func(l *lua.LState) int {
		args := make(map[string]interface{})
		str := l.ToString(1)
		err := json.Unmarshal([]byte(str), &args)
		if err != nil {
			log.Println("Error: Json Decode fail\n" + string(str))
			l.Push(lua.LString("ERR"))
			return 1
		}
		if err := this.Db.ModifyRow(notnull, args); err != nil {
			log.Println(err)
			l.Push(lua.LString("FAIL"))
			return 1
		} else {
			l.Push(lua.LString("SUCC"))
			return 1
		}
	}
}

func (this VMs) DelRowFunc() func(l *lua.LState) int {
	return func(l *lua.LState) int {
		var id int
		if str, ok := l.Get(-1).(lua.LString); ok {
			l.Pop(1)
			args := make(map[string]interface{})
			err := json.Unmarshal([]byte(str), &args)
			id, _ = strconv.Atoi(args[":id"].(string))
			if err != nil {
				log.Println("Error: Json Decode fail\n" + string(str))
				l.Push(lua.LString("ERR"))
				return 1
			}

			if err != nil {
				log.Println("ERROR: Wrong Input\n\t" + string(id))
				return 0
			}
			if err = this.Db.DeleteRow(id); err != nil {
				l.Push(lua.LString("SUCC"))
			} else {
				l.Push(lua.LString("FAIL"))
			}
		} else {
			l.Push(lua.LString("ERR"))
		}
		return 1
	}
}

func (this VMs) FindRowFunc() func(l *lua.LState) int {
	return func(l *lua.LState) int {
		if str, ok := l.Get(-1).(lua.LString); ok {
			args := make(map[string]interface{})
			err := json.Unmarshal([]byte(str), &args)
			id, _ := args[":id"].(string)

			Id, err := strconv.Atoi(id)
			if err != nil {
				log.Println(err)
				log.Println(string(id))
				l.Push(lua.LString("ERR"))
				return 1
			}
			ret, err := this.Db.FindRow(Id)
			if err != nil {
				log.Println(err)
				l.Push(lua.LString("ERR"))
				return 1
			}
			l.Push(lua.LString(ret))
			return 1
		} else {
			l.Push(lua.LString("ERR"))
			return 1
		}
	}
}
