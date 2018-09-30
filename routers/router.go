package routers

import (
	"mztw/controllers"
	"github.com/astaxie/beego"
)

func init() {
    beego.Router("/", &controllers.MainController{})
	beego.Router("/getgradbta", &controllers.GrabdataDisat{}, "Get:Getgrabdata")
}
