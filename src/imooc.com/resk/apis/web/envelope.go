package web

import (
	"github.com/kataras/iris"
	"resk-projects/infra"
	"resk-projects/infra/base"
	"resk-projects/services"
)

func init() {
	infra.RegisterApi(new(EnvelopeApi))
}

type EnvelopeApi struct {
	service services.RedEnvelopeService
}


func (e *EnvelopeApi) Init() {
	e.service = services.GetRedEnvelopeService()
	groupRouter := base.Iris().Party("/v1/envelope")
	groupRouter.Post("/sendout", e.sendOutHandler)
	groupRouter.Post("/receive", e.receiveHandler)
}

//收红包
func (e *EnvelopeApi) receiveHandler(ctx iris.Context)  {
	dto := services.RedEnvelopeReceiveDTO{}
	err := ctx.ReadJSON(&dto)
	r := base.Res{
		Code: base.ResCodeOk,
	}
	if err != nil {
		r.Code = base.ResCodeRequestParamsError
		r.Message = err.Error()
		ctx.JSON(&r)
		return
	}
	item, err := e.service.Receive(dto)
	if err != nil {
		r.Code = base.ResCodeInnerServerError
		r.Message = err.Error()
		ctx.JSON(&r)
		return
	}
	r.Data = item
	ctx.JSON(&r)
}
func (e *EnvelopeApi) sendOutHandler(ctx iris.Context) {
	dto := services.RedEnvelopeSendingDTO{}
	err := ctx.ReadJSON(&dto)
	r := base.Res{
		Code: base.ResCodeOk,
	}
	if err != nil {
		r.Code = base.ResCodeRequestParamsError
		r.Message = err.Error()
		ctx.JSON(&r)
		return
	}
	activity, err := e.service.SendOut(dto)
	if err != nil {
		r.Code = base.ResCodeInnerServerError
		r.Message = err.Error()
		ctx.JSON(&r)
		return
	}
	r.Data = activity
	ctx.JSON(&r)
}
