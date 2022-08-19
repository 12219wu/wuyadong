package routers

import (
	beego "github.com/beego/beego/v2/server/web"
	"netty_jtt808/controllers"
)

func init() {
	beego.Router("/", &controllers.MainController{})
}
