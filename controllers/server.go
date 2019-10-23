package controllers

import (
	"net/http"

	"github.com/astaxie/beego"
	"github.com/gorilla/websocket"
)

// ServerController ServerController
type ServerController struct {
	beego.Controller
}

// Get Get Info
func (c *ServerController) Get() {
	name := c.GetString("name")
	if len(name) == 0 {
		beego.Error("name is null")
		c.Redirect("/", 302)
		return
	}
	beego.Info("get name:" + name + ", and send to socketChat.html")
	c.Data["name"] = name
	c.TplName = "socketChat.html"
}

// WebSocket WebSocket连接
func (c *ServerController) WebSocket() {
	// 相当于验证参数 name
	name := c.GetString("name")
	if len(name) == 0 {
		beego.Error("name is null")
		c.Redirect("/", 302)
		return
	}

	// 检测http头中upgrader属性,若为websocket,则将http协议升级
	conn, err := (&websocket.Upgrader{}).Upgrade(c.Ctx.ResponseWriter, c.Ctx.Request, nil)

	if _, ok := err.(websocket.HandshakeError); ok {
		beego.Error("Not a websocket connection")
		http.Error(c.Ctx.ResponseWriter, "not a wensocket handshake", 400)
		return
	} else if err != nil {
		beego.Error("Cannot setup Websockt connection:", err)
		return
	}

	var client Client
	client.conn = conn
	client.name = name

	// 如果用户列表中没有该用户
	if !clients[client] {
		join <- client
		beego.Info("user:", client.name, " websocket connect success")
	}

	// 当函数返回时，将该用户加入退出通道，并断开用户连接
	defer func() {
		leave <- client
		client.conn.Close()
	}()

	// 因为websocket是长连接,一次握手,创建持久性的连接 直到连接断开
	for {
		// 读取消息 如果断开连接 则会返回错误
		_, msgStr, err := client.conn.ReadMessage()
		// 如果返回错误 就退出循环
		if err != nil {
			break
		}

		beego.Info("websocket------------receive:" + string(msgStr))

		// 如果没有错误,则把用户发送的消息放入message通道中
		var msg Message
		msg.Name = client.name
		msg.EventType = 0
		msg.Message = string(msgStr)
		message <- msg
	}
}
