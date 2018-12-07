package controllers

import (
	"fmt"
	"github.com/astaxie/beego"
	"github.com/du2016/web-terminal-in-go/container-webshell/models"
	"github.com/gorilla/websocket"
	"log"
	"net"
)

type Wscontroller struct {
	beego.Controller
}

var upgrader = websocket.Upgrader{}

func (self *Wscontroller) Get() {
	host := self.Input().Get("h")
	port := self.Input().Get("p")
	containerid := self.Input().Get(("containers_id"))
	rows := self.Input().Get("rows")
	cols := self.Input().Get("cols")
	execid := models.Getexecid(host, port, containerid)
	log.Println("execid is", execid)
	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%s", host, port))
	if err != nil {
		beego.Error(err)
	}
	data := "{\"Tty\":true}"
	_, err = conn.Write([]byte(fmt.Sprintf("POST /exec/%s/start HTTP/1.1\r\nHost: %s\r\nContent-Type: application/json\r\nContent-Length: %s\r\n\r\n%s", execid, fmt.Sprintf("%s:%s", host, port), fmt.Sprint(len([]byte(data))), data)))
	if err != nil {
		log.Println(err)
	}
	models.Resizecontainer(host, port, execid, rows, cols)
	ws, err := upgrader.Upgrade(self.Ctx.ResponseWriter, self.Ctx.Request, nil)
	c := &models.Connection{Send: make(chan []byte, 256), Ws: ws}
	go c.Writer(conn)
	c.Reader(conn)
	defer conn.Close()
}
