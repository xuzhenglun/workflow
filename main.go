package main

import (
	"github.com/xuzhenglun/workflow/api"
	"github.com/xuzhenglun/workflow/core"
	"github.com/xuzhenglun/workflow/restful"
)

func main() {
	core := core.InitCore()
	core.RegeditApi(api.List)

	restful.Run(8080, core)
}
