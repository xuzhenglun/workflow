package controllers

import (
	"log"

	"github.com/astaxie/beego"
	"github.com/xuzhenglun/workflow/front/models"
)

type DefaultController struct {
	beego.Controller
}

func (c *DefaultController) Get() {
	auth := c.Ctx.GetCookie("auth")
	log.Println(auth)
	if auth == "" {
		c.Redirect("/login", 302)
	}

	c.Layout = "farmework.tpl"
	c.TplName = "default.tpl"
	activity := c.GetString(":activity")
	id := c.GetString(":id")

	info := model.GetInfo(beego.AppConfig.String("url")+"/"+activity+"_reader/"+id, auth)
	if info == nil {
		log.Println("Can't find reader Activity")
		c.Redirect("/", 302)
	}

	var showInfo map[string]string
	switch info.Code {
	case 200:
		showInfo = info.Msg
	}
	log.Println(info, showInfo)

	args := model.GetArgs(beego.AppConfig.String("url")+"/"+activity+"/help", auth)
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
	log.Println(id, args, activity)
	c.Data["Args"] = needargs
	c.Data["Name"] = activity
	c.Data["IsPass"] = pass
	c.Data["Info"] = showInfo
	c.Data["Activities"] = a
	c.Data["User"] = model.GetUserInfo(auth)
}

func (this *DefaultController) Post() {
	this.EnableRender = false

	resp, err := model.Post(beego.AppConfig.String("url")+this.Ctx.Request.RequestURI, this.Ctx.Request)

	if err != nil {
		this.Ctx.ResponseWriter.Write(resp)
		return
	}
	this.Ctx.ResponseWriter.Write(resp)
}
