package core

import (
	"fmt"
	"log"
	"time"

	"github.com/yuin/gluamapper"
	"github.com/yuin/gopher-lua"
)

type Activity struct {
	L        *lua.LState
	Father   string
	Pass     string
	Name     string
	Types    string
	NeedArgs string
	Handler  lua.LValue
}

type HandleErr struct {
	When time.Time
	What string
}

func FindActivityByName(l *lua.LState, name string) *Activity {
	activity := Activity{L: l}

	if table, ok := l.GetGlobal(name).(*lua.LTable); ok {
		if err := gluamapper.Map(table, &activity); err != nil {
			log.Println(err)
			log.Println(`Waring: Faild to find the activity "` + name + `"`)
			return nil
		} else {
			if activity.Pass == "" || activity.Pass == "false" {
				activity.Pass = "0"
			} else if activity.Pass == "true" {
				activity.Pass = "1"
			}
			return &activity
		}
	} else {
		log.Println("Can't find \"" + name + "\"")
		return nil
	}
}

func (activity *Activity) Handle(param string) (string, error) {
	if activity == nil {
		panic("Nil Activity")
	}

	log.Println("Handle Args: " + param)
	l := (*activity).L
	err := l.CallByParam(lua.P{
		Fn:      (*activity).Handler,
		NRet:    1,
		Protect: true,
	}, lua.LString(param))

	if err != nil {
		log.Println(err)
		return "", HandleErr{When: time.Now(), What: "VM Runtime Error"}
	}

	lret := l.Get(-1)
	if ret, ok := lret.(lua.LString); ok {
		defer l.Pop(1)
		return string(ret), nil
	} else {
		log.Println("Warring:can't get any return")
		return "", HandleErr{When: time.Now(), What: "Return type should be string"}
	}
}

func (e HandleErr) Error() string {
	return fmt.Sprintf("%v: %v", e.When, e.What)
}
