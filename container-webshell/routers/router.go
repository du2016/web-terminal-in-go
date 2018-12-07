package routers

import (
	"github.com/astaxie/beego"
	"github.com/du2016/web-terminal-in-go/container-webshell/controllers"
)

func init() {
	beego.Router("/", &controllers.MainController{})
	beego.Router("/ws", &controllers.Wscontroller{})
}
