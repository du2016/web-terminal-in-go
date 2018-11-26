package routers

import (
	"github.com/du2016/web-terminal-in-go/k8s-webshell/controllers"
	"github.com/astaxie/beego"
	"golang.org/x/net/websocket"
)

func init() {
    	beego.Router("/", &controllers.HomeController{})
	beego.Router("/terminal", &controllers.TerminalController{}, "get:Get")
	beego.Handler("/terminal/ws", websocket.Handler(controllers.EchoHandler))
}
