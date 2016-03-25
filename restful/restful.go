package restful

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/drone/routes"
	"github.com/xuzhenglun/workflow/core"
)

func Run(port int, vms core.CoreIoBus) {
	mux := routes.New()

	mux.Get("/:activity/:id", HandlerHub(vms))
	mux.Post("/:activity/:id", HandlerHub(vms))

	http.ListenAndServe(":"+strconv.Itoa(port), mux)
}

func HandlerHub(vms core.CoreIoBus) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		param := r.URL.Query()
		activity := param.Get(":activity")
		log.Println(activity)
		log.Println(vms.GetMapper()[activity], ":", r.Method)
		if strings.Contains(vms.GetMapper()[activity], r.Method) == false {
			fmt.Fprintln(w, "Wrong way to access. Use \""+vms.GetMapper()[activity]+"\" Please")
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
