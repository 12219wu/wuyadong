package presenters

import (
	"JTTServer/jtt"
	"log"
)

type LoginPresenter struct {
	jtt.BasePresenter
}

func (l *LoginPresenter) TerminalAuth() {
	if msg, ok := l.Ctx.Message().(*jtt.MsgTerminalAuth); ok {
		log.Printf("%s->%s 终端鉴权 %v", l.Ctx.Client().RemoteAddr(), l.Ctx.Client().LocalAddr(), msg)
		resp := jtt.NewMsgServerResponse(msg.Number, msg.ID, 0)
		l.Ctx.Response(resp)
	}
}

func (l *LoginPresenter) PositionReport() {
	if msg, ok := l.Ctx.Message().(*jtt.MsgPositionReport); ok {
		log.Printf("%s->%s 终端位置上报 %v", l.Ctx.Client().RemoteAddr(), l.Ctx.Client().LocalAddr(), msg)
		resp := jtt.NewMsgServerResponse(msg.Number, msg.ID, 0)
		l.Ctx.Response(resp)
	}
}

func (l *LoginPresenter) PositionBatchReport() {
	if msg, ok := l.Ctx.Message().(*jtt.MsgPosBatchReport); ok {
		log.Printf("%s->%s 终端位置批量上报 %v", l.Ctx.Client().RemoteAddr(), l.Ctx.Client().LocalAddr(), msg)
		resp := jtt.NewMsgServerResponse(msg.Number, msg.ID, 0)
		l.Ctx.Response(resp)
	}
}

func (l *LoginPresenter) TerminalHeatbeat() {
	if msg, ok := l.Ctx.Message().(*jtt.MsgTerminalHeartbeat); ok {
		log.Printf("%s->%s 终端心跳，%v", l.Ctx.Client().RemoteAddr(), l.Ctx.Client().LocalAddr(), msg)
		resp := jtt.NewMsgServerResponse(msg.Number, msg.ID, 0)
		l.Ctx.Response(resp)
	}
}
