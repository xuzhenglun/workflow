package controllers

import (
	"log"

	"github.com/astaxie/beego"
	"github.com/xuzhenglun/workflow/front/models"
)

type StartController struct {
	beego.Controller
}

func (this *StartController) Get() {
	auth := this.Ctx.GetCookie("auth")

	this.Layout = "farmework.tpl"
	this.TplName = "start.tpl"

	args := model.GetArgs(beego.AppConfig.String("url")+"/start/help", auth)
	var needargs []string

	pass := false
	for _, v := range args.Args {
		if v == "pass" {
			pass = true
		} else {
			needargs = append(needargs, v)
		}
	}

	a := model.GetAllActivites(beego.AppConfig.String("url"), auth)
	this.Data["Args"] = needargs
	this.Data["Name"] = "填报申请"
	this.Data["IsPass"] = pass
	this.Data["Activities"] = a
	this.Data["User"] = model.GetUserInfo(auth)
}

func (this *StartController) Post() {
	this.EnableRender = false
	log.Println(this.Ctx.Request.RequestURI)
	resp, err := model.Post(beego.AppConfig.String("url")+this.Ctx.Request.RequestURI+"/new", this.Ctx.Request)

	if err != nil {
		this.Ctx.ResponseWriter.Write(resp)
		return
	}
	this.Ctx.ResponseWriter.Write(resp)

}
