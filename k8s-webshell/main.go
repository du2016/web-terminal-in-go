package main

import (
	_ "github.com/du2016/web-terminal-in-go/k8s-webshell/routers"
	"github.com/astaxie/beego"
)

func main() {
	beego.Run()
}

