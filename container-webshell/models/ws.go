package models

import (
	"github.com/astaxie/beego"
	"github.com/gorilla/websocket"
	"net"
	//"fmt"
	"encoding/json"
	"log"
)

type Connection struct {
	// websocket 连接器
	Ws *websocket.Conn

	// 发送信息的缓冲 channel
	Send chan []byte
}

func (c *Connection) Reader(conn net.Conn) {

	for {
		_, message, err := c.Ws.ReadMessage()
		if err != nil {
			log.Println(err)
		}
		conn.Write(message)
	}
}

type tojson struct {
	Data string `json:"data"`
}

func (c *Connection) Writer(conn net.Conn) {
	//c.Ws.Close()
	for {
		b := make([]byte, 512)
		conn.Read(b)
		j := tojson{Data: string(b)}
		d, _ := json.Marshal(j)
		c.Ws.WriteMessage(websocket.TextMessage, d)
		beego.Error(string(b))
	}
}
