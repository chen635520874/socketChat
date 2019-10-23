package main

import (
	_ "mygoproject/beego/socketchat/routers"
	"github.com/astaxie/beego"
)

func main() {
	beego.Run()
}

