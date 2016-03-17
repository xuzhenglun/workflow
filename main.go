package main

import (
	"github.com/xuzhenglun/workflow/api"
	"github.com/xuzhenglun/workflow/cli"
	"github.com/xuzhenglun/workflow/core"
	"github.com/xuzhenglun/workflow/database"
	"github.com/xuzhenglun/workflow/restful"
)

func main() {
	core := core.InitCore()
	core.RegeditApi(api.List)
	db := database.NewMysql("root:xuzl93042136@/workflow?charset=utf8")
	core.SetDataBase(db)

	go cli.Run(core)
	restful.Run(8080, core)
}
