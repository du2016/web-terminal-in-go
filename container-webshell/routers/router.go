package routers

import (
	"github.com/du2016/web-terminal-in-go/container-webshell/controllers"
	"github.com/astaxie/beego"
)

func init() {
    beego.Router("/", &controllers.MainController{})
	beego.Router("/ws", &controllers.Wscontroller{})
}
