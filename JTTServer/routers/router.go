package routers

import (
	"JTTServer/controllers"
	"JTTServer/jtt"
	"JTTServer/presenters"

	beego "github.com/beego/beego/v2/server/web"
)

func init() {
	beego.Router("/", &controllers.MainController{})

	jtt.Router(jtt.MsgIDTerminalAuth, &presenters.LoginPresenter{}, "TerminalAuth")
	jtt.Router(jtt.MsgIDPositionReport, &presenters.LoginPresenter{}, "PositionReport")
	jtt.Router(jtt.MsgIDPositionBatchReport, &presenters.LoginPresenter{}, "PositionBatchReport")
	jtt.Router(jtt.MsgIDTerminalHeartbeat, &presenters.LoginPresenter{}, "TerminalHeatbeat")
}
