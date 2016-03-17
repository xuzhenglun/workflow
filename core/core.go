package core

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"regexp"
	"time"

	"github.com/layeh/gopher-json"
	"github.com/yuin/gopher-lua"
)

const SCRIPT_DIR = "./scripts/"

type VMs struct {
	Scripts    string
	Activities map[string]string
	Mappler    map[string]string
	Api        map[string]func(*lua.LState) int
	Db         DataBase
	Re         *regexp.Regexp
}

type VM struct {
	lua.LState
}

func InitCore() *VMs {
	var this VMs
	path, _ := os.Open(SCRIPT_DIR)
	info, _ := path.Readdir(-1)
	re, _ := regexp.Compile(`^\w+[^/.]`)
	this.Activities = make(map[string]string)
	for _, v := range info {
		f, err := ioutil.ReadFile(SCRIPT_DIR + v.Name())
		if err != nil {
			log.Panic(err)
		}
		this.Scripts = this.Scripts + string(f)
		name := re.Find([]byte(v.Name()))
		this.Activities[string(name)] = string(f)
		log.Println("Script \"" + string(name) + "\" loaded")
	}
	this.InitMap()
	this.Re, _ = regexp.Compile(`\w+`)
	return &this
}

func (this *VMs) RequestHandler(w ReponseWriter, r *Request) error {
	l := lua.NewState()
	defer l.Close()

	if this.Api != nil {
		for key, value := range this.Api {
			l.SetGlobal(key, l.NewFunction(value))
		}
	}
	json.Preload(l)

	if v, ok := this.Activities[r.Name]; ok {
		l.DoString(v)
	} else {
		log.Println(`Warning: Failed to Find "` + r.Name + `.lua", I'will search it form globle, It's may cause performance issue.`)
		l.DoString(this.Scripts)
	}

	activity := FindActivityByName(l, r.Name)
	if activity == nil {
		return HandleErr{When: time.Now(), What: "Can't find activity"}
	}

	this.InitDatabase(l, activity)
	ret, err := activity.Handle(r.Args)
	if err != nil {
		log.Println(err)
		return err
	}

	fmt.Fprintln(w, ret)
	return err
}

func (this *VMs) InitMap() {
	this.Mappler = make(map[string]string)
	l := lua.NewState()
	defer l.Close()
	l.DoString(this.Scripts)
	for i, _ := range this.Activities {
		if activity, ok := l.GetGlobal(i).(*lua.LTable); ok {
			if value, ok := activity.RawGetString("typ").(lua.LString); ok {
				this.Mappler[i] = string(value)
			} else {
				log.Println("Error: Can't find field \"typ\" in Table \"" + i + "\"")
			}
		} else {
			log.Println("Error: Can't find Table \"" + i + "\"")
		}
	}
}

func (this *VMs) GetMapper() map[string]string {
	return this.Mappler
}

func (this *VMs) RegeditApi(list map[string]func(*lua.LState) int) {
	if list != nil {
		this.Api = list
	}
}

func (this *VMs) SetDataBase(database DataBase) {
	this.Db = database
}

func (this *VMs) RawLuaHandler(str string) string {
	l := lua.NewState()
	if err := l.DoString(str); err != nil {
		return err.Error()
	}
	return "Done"

}
