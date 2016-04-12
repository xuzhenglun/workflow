package controllers

import (
	"github.com/astaxie/beego"
	"github.com/xuzhenglun/workflow/front/models"
)

type DashboardController struct {
	beego.Controller
}

func (c *DashboardController) Get() {
	c.Layout = "farmework.tpl"
	c.TplName = "dashboard.tpl"
	auth := c.Ctx.GetCookie("auth")
	if auth == "" {
		c.Redirect("/login", 302)
		return
	}

	a := model.GetAllActivites(beego.AppConfig.String("url"), auth)

	c.Data["Activities"] = a
	c.Data["User"] = model.GetUserInfo(auth)

	name := c.GetString(":activity")
	if name == "dashboard" {
		c.Data["Name"] = ""
		return
	}

	ids := model.GetAllActivites(beego.AppConfig.String("url")+"/"+name, auth)
	c.Data["Name"] = name
	c.Data["ids"] = ids
}
