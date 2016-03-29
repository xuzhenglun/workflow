package main

import (
	"github.com/xuzhenglun/workflow/api"
	"github.com/xuzhenglun/workflow/cli"
	"github.com/xuzhenglun/workflow/core"
	"github.com/xuzhenglun/workflow/database"
	"github.com/xuzhenglun/workflow/restful"
	"github.com/xuzhenglun/workflow/signature"
)

func main() {
	core := core.InitCore()
	core.RegeditApi(api.List)

	db := database.NewMongoDB("")
	//db := database.NewMysql("user:passwd@/workflow?charset=utf8")
	core.SetDataBase(db)

	sign := signature.NewSigner("keys.json", db)
	core.SetAuther(sign)

	go cli.Run(core)
	restful.Run(8080, core)
}
