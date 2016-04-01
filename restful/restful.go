package restful

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/bitly/go-simplejson"
	"github.com/drone/routes"
	"github.com/xuzhenglun/workflow/core"
)

func Run(port int, vms core.CoreIoBus) {
	mux := routes.New()

	mux.Get("/", GetPurviewActivities(vms))
	mux.Get("/:activity/help", listArgs(vms))
	mux.Get("/:activity/:id", HandlerHub(vms))
	mux.Post("/:activity/:id", HandlerHub(vms))
	mux.Get("/:action", ListActivites(vms))

	http.ListenAndServe(":"+strconv.Itoa(port), mux)
}

func HandlerHub(vms core.CoreIoBus) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		auth := r.Header.Get("auth")

		param := r.URL.Query()
		activity := param.Get(":activity")
		log.Println(activity, ":", vms.GetMapper()[activity], ":", r.Method)
		if strings.Contains(vms.GetMapper()[activity], r.Method) == false {
			w.Write([]byte("Wrong way to access. Use \"" + vms.GetMapper()[activity] + "\" Please"))
			return
		}
		id := param.Get(":id")
		args := `{":id":"` + id + `"}`

		if r.Method == "POST" {
			if err := r.ParseForm(); err != nil {
				log.Println(err)
			}

			tmp := make(map[string]string)
			for k, v := range r.Form {
				tmp[k] = v[0]
			}

			a, _ := json.Marshal(tmp)
			args = string(a)
		}

		requestVM := core.Request{
			Name: activity,
			Args: args,
			Id:   id,
			Auth: []byte(auth),
		}

		responseVM := core.Response{}
		if err := vms.RequestHandler(&responseVM, &requestVM); err != nil {
			log.Println("Request Handler Error")
			fmt.Fprint(w, err)
			w.WriteHeader(500)
			return
		} else {
			if responseVM.Status == -1 {
				log.Println("VM Return Error Signal -1")
				w.WriteHeader(404)
				return
			} else {
				fmt.Fprint(w, responseVM.Body)
				return
			}
		}
	}
}

func ListActivites(vms core.CoreIoBus) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		auth := r.Header.Get("auth")

		param := r.URL.Query()
		action := param.Get(":action")
		list, err := vms.GetActivities([]byte(auth), action)
		if err == nil {
			fmt.Fprint(w, list)
		} else {
			w.WriteHeader(404)
		}
	}
}
func listArgs(vms core.CoreIoBus) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		auth := r.Header.Get("auth")

		param := r.URL.Query()
		activity := param.Get(":activity")

		ret := simplejson.New()
		if l, err := vms.ListNeedArgs([]byte(auth), activity); err == nil {
			ret.Set("args", l)
		} else {
			log.Println(err)
			return
		}
		if c, err := ret.Encode(); err == nil {
			w.Write(c)
		} else {
			log.Println(err)
		}
	}
}

func GetPurviewActivities(vms core.CoreIoBus) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		auth := r.Header.Get("auth")

		resp, err := vms.GetPurviewActivities(auth)
		if err != nil {
			log.Println(err)
		}

		ret, err := json.Marshal(resp)

		if err != nil {
			log.Println(err)
			w.WriteHeader(403)
		} else {
			w.Write(ret)
		}
	}
}
