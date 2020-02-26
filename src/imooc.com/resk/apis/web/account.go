package web

import (
	"github.com/kataras/iris"
	"resk-projects/infra/base"
	"resk-projects/services"
)

//定义web api的时候，对每一个子业务定义统一的前缀
//资金账户的根目录定义为：/account
//版本号: /v1/account
var groupRouter iris.Party

func init() {
	groupRouter = base.Iris().Party("/v1/account")
	create(groupRouter)
}

//账户创建接口: /v1/account/create
//POST body json
func create(groupRouter iris.Party) {
	groupRouter.Post("/create",
		func(ctx iris.Context) {
			//获取请求参数
			account :=services.AccountCreatedDTO{}
			err := ctx.ReadJSON(&account)
			r:=base.Res{
				Code:base.ResCodeOk,
			}
			if err!=nil{
				r.Code=base.ResCodeRequestParamsError
				r.Message=err.Error()
				ctx.JSON(&r)
				return
			}
			//执行创建账户的代码
		})
}

//转账接口: /v1/account/transfer

//查询红包账户接口：/v1/account/envelope/get

//查询账户信息接口：/v1/account/get
