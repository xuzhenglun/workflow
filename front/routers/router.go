package routers

import "github.com/astaxie/beego"
import "github.com/xuzhenglun/workflow/front/controllers"

func init() {
	beego.Router("/", &controllers.IndexController{})
	beego.Router("/start", &controllers.StartController{})
	beego.Router("/logout", &controllers.LogOutController{})
	beego.Router("/login", &controllers.LoginController{})

	beego.Router("/:activity", &controllers.DashboardController{})
	beego.Router("/:activity/:id", &controllers.DefaultController{})
	beego.SetStaticPath("/static", "static")
}
