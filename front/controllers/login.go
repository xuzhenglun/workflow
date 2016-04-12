package controllers

import (
	"github.com/astaxie/beego"
	"github.com/xuzhenglun/workflow/front/models"
)

type LoginController struct {
	beego.Controller
}

func (this *LoginController) Get() {
	this.Layout = "farmework.tpl"
	this.TplName = "login.tpl"
	auth := this.Ctx.GetCookie("auth")
	if auth != "" {
		if a := model.GetAllActivites(beego.AppConfig.String("url"), auth); a != nil {
			this.Redirect("/", 302)
		}
	}
}

func (this *LoginController) Post() {
	this.Layout = "farmework.tpl"
	this.TplName = "login.tpl"
	license := this.Input().Get("auth")
	long := this.Input().Get("long") == "on"

	maxAge := 0
	if long {
		maxAge = 1<<31 - 1
	}

	if license == "" {
		this.Redirect("/login", 302)
	}
	if a := model.GetAllActivites(beego.AppConfig.String("url"), license); a != nil {
		this.Ctx.SetCookie("auth", license, maxAge)
		this.Redirect("/", 302)
	} else {
		this.Redirect("/login", 302)
	}
}

type LogOutController struct {
	beego.Controller
}

func (this *LogOutController) Get() {
	this.EnableRender = false
	this.Ctx.SetCookie("auth", "", -1)
	this.Redirect("/", 302)
}
