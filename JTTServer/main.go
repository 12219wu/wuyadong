package main

import (
	"JTTServer/jtt"
	_ "JTTServer/routers"

	beego "github.com/beego/beego/v2/server/web"
)

func main() {
	jtt.Run("127.0.0.1:8081")
	beego.Run()
}
