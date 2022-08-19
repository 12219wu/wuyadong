package main

import (
	beego "github.com/beego/beego/v2/server/web"
	_ "netty_jtt808/routers"
)

func main() {
	beego.Run()
}
