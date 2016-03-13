package restful

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/drone/routes"
	"github.com/xuzhenglun/workflow/core"
)

func Run(port int, vms core.CoreIoBus) {
	mux := routes.New()

	mux.Get("/:activity/:args", HandlerHub(vms))
	mux.Post("/:activity/:args", HandlerHub(vms))

	http.ListenAndServe(":8080", mux)
}

func HandlerHub(vms core.CoreIoBus) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		param := r.URL.Query()
		activity := param.Get(":activity")
		if strings.Contains(r.Method, vms.GetMapper()[activity]) == false {
			fmt.Sprintln(w, "Wrong way to access. Use \""+r.Method+"\" Please")
		}
		args := param.Get(":args")
		requestVM := core.Request{
			Name: activity,
			Args: args,
		}

		responseVM := core.Response{}
		if err := vms.RequestHandler(&responseVM, &requestVM); err != nil {
			log.Println("Request Handler Error")
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
