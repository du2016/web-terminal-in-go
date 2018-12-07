package main

import (
	"github.com/astaxie/beego"
	_ "github.com/du2016/web-terminal-in-go/k8s-webshell/routers"
)

func main() {
	beego.Run()
}
