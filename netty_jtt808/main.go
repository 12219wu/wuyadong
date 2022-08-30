package main

import (
	beego "github.com/beego/beego/v2/server/web"
	"netty_jtt808/client"
	_ "netty_jtt808/routers"
)

func main() {
	client.SetUpCodec()
	beego.Run()
}
