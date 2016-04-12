package controllers

import (
	"github.com/astaxie/beego"
	"github.com/xuzhenglun/workflow/front/models"
)

type IndexController struct {
	beego.Controller
}

func (this *IndexController) Get() {
	auth := this.Ctx.GetCookie("auth")
	this.TplName = "index.tpl"
	this.Data["User"] = model.GetUserInfo(auth)
}
