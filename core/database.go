package core

import (
	"log"
	"regexp"

	"github.com/bitly/go-simplejson"
	"github.com/yuin/gopher-lua"
)

func (this VMs) InitDatabase(l *lua.LState, a *Activity) {
	l.SetGlobal("AddRow", l.NewFunction(this.AddRowFunc))
	l.SetGlobal("DelRow", l.NewFunction(this.DelRowFunc))
	l.SetGlobal("ModRow", l.NewFunction(this.ModRowFunc))
	l.SetGlobal("FindRow", l.NewFunction(this.FindRowFunc))
}

func NewResult(code int, r interface{}) lua.LString {
	j := simplejson.New()
	j.Set("Code", code)
	if v, ok := r.(error); ok {
		j.Set("Msg", v.Error())
	}
	if v, ok := r.(string); ok {
		j.Set("Msg", v)
	}

	ret, err := j.MarshalJSON()
	if err == nil {
		return lua.LString(ret)
	} else {
		log.Println(err)
		return lua.LString(`{"Code":500,"Msg":"UnExpect Error"}`)
	}
}

var re, _ = regexp.Compile(`\w+`)

func Splite(str lua.LString) []string {
	return re.FindAllString(string(str), -1)
}

func (this VMs) AddRowFunc(l *lua.LState) int {
	env := l.ToTable(1)
	req, err := simplejson.NewJson([]byte(l.ToString(2)))
	if err != nil {
		l.Push(NewResult(500, err))
		return 1
	}

	var needArgs []string
	if v, ok := env.RawGetString("needArgs").(lua.LString); ok {
		needArgs = Splite(v)
	}

	events := make(map[string]string)
	process := make(map[string]string)
	process["JustDone"] = env.RawGetString("name").String()

	for _, k := range needArgs {
		v, err := req.Get(k).String()
		if err != nil {
			log.Println(err)
			l.Push(NewResult(401, err))
			return 1
		} else {
			if k != "Pass" || k != "Done" {
				events[k] = v
			} else {
				process[k] = v
			}
		}
	}

	//Todo: Maybe I can add User Group here
	log.Println(events, process)
	if err := this.Db.AddRow(events, process); err != nil {
		log.Println(err)
		l.Push(NewResult(500, err))
	} else {
		l.Push(NewResult(200, "Sucess"))
	}

	return 1
}

func (this VMs) ModRowFunc(l *lua.LState) int {
	env := l.ToTable(1)
	req, err := simplejson.NewJson([]byte(l.ToString(2)))
	if err != nil {
		l.Push(NewResult(500, err))
		return 1
	}

	var needArgs []string
	if v, ok := env.RawGetString("needArgs").(lua.LString); ok {
		needArgs = Splite(v)
	}

	events := make(map[string]string)
	process := make(map[string]string)
	process["JustDone"] = env.RawGetString("name").String()

	for _, k := range needArgs {
		v, err := req.Get(k).String()
		if err != nil {
			log.Println(err)
			l.Push(NewResult(401, err))
			return 1
		} else {
			if k != "Pass" || k != "Done" {
				events[k] = v
			} else {
				process[k] = v
			}
		}
	}

	id, err := req.Get(":id").String()
	if err == nil {
		process["id"] = id
	} else {
		l.Push(NewResult(401, err))
		return 1
	}

	//Todo: Maybe I can add User Group here
	if err := this.Db.ModifyRow(events, process); err != nil {
		log.Println(err)
		l.Push(NewResult(500, err))
	} else {
		l.Push(NewResult(200, "Sucess"))
	}

	return 1
}

func (this VMs) DelRowFunc(l *lua.LState) int {
	//env := l.ToTable(1)
	req, err := simplejson.NewJson([]byte(l.ToString(2)))
	if err != nil {
		l.Push(NewResult(500, err))
		return 1
	}

	id, err := req.Get(":id").String()
	if err != nil {
		log.Println(err)
		l.Push(NewResult(401, err))
	} else {
		err = this.Db.DeleteRow(id)
		if err != nil {
			l.Push(NewResult(401, err))
		} else {
			l.Push(NewResult(200, "Sucess"))
		}
	}

	return 1
}

func (this VMs) FindRowFunc(l *lua.LState) int {
	//env := l.ToTable(1)
	req, err := simplejson.NewJson([]byte(l.ToString(2)))
	if err != nil {
		l.Push(NewResult(500, err))
	}

	var res string
	id, err := req.Get(":id").String()
	if err != nil {
		log.Println(err)
		l.Push(NewResult(401, err))
	} else {
		res, err = this.Db.FindRow(id)
		if err != nil {
			l.Push(NewResult(401, err))
		} else {
			l.Push(NewResult(200, res))
		}
	}

	return 1
}
