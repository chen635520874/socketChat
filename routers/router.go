package routers

import (
	"mygoproject/beego/socketchat/controllers"

	"github.com/astaxie/beego"
)

func init() {
	beego.Router("/", &controllers.MainController{})
	beego.Router("/socketChat", &controllers.ServerController{})
	beego.Router("/socketChat/WS", &controllers.ServerController{}, "get:WebSocket")
}
