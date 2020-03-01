package web

import (
	"github.com/kataras/iris"
	"resk-projects/infra"
	"resk-projects/infra/base"
	"resk-projects/services"
)

//定义web api的时候，对每一个子业务定义统一的前缀
//资金账户的根目录定义为：/account
//版本号: /v1/account
const (
	ResCodeBizTransferFailure = base.ResCode(610)
)

func init() {
	infra.RegisterApi(new(AccountApi))
}

type AccountApi struct {
}

func (a *AccountApi) Init() {
	groupRouter := base.Iris().Party("/v1/account")
	groupRouter.Post("/create", createHandler)
	groupRouter.Post("/transfer", transferHandler)
	groupRouter.Get("/envelope/get", getEnvelopeAccountHandler)
	groupRouter.Get("/get", getAccountHandler)
}

//账户创建接口: /v1/account/createHandler
//POST body json
func createHandler(ctx iris.Context) {
	//获取请求参数
	account := services.AccountCreatedDTO{}
	err := ctx.ReadJSON(&account)
	r := base.Res{
		Code: base.ResCodeOk,
	}
	if err != nil {
		r.Code = base.ResCodeRequestParamsError
		r.Message = err.Error()
		ctx.JSON(&r)
		return
	}
	//执行创建账户的代码
	service := services.GetAccountService()
	dto, err := service.CreateAccount(account)
	if err != nil {
		r.Code = base.ResCodeInnerServerError
		r.Message = err.Error()
	}
	r.Data = dto
	ctx.JSON(&r)
}

//转账接口: /v1/account/transfer
func transferHandler(ctx iris.Context) {
	//获取请求参数
	account := services.AccountTransferDTO{}
	err := ctx.ReadJSON(&account)
	r := base.Res{
		Code: base.ResCodeOk,
	}
	if err != nil {
		r.Code = base.ResCodeRequestParamsError
		r.Message = err.Error()
		ctx.JSON(&r)
		return
	}
	//执行转账逻辑
	service := services.GetAccountService()
	status, err := service.Transfer(account)
	if err != nil {
		r.Code = base.ResCodeInnerServerError
		r.Message = err.Error()
	}
	r.Data = status
	if status != services.TransferStatusSuccess {
		r.Code = ResCodeBizTransferFailure
		r.Message = err.Error()
	}
	ctx.JSON(&r)

}

//查询红包账户接口：/v1/account/envelope/get
func getEnvelopeAccountHandler(ctx iris.Context) {
	userId := ctx.URLParam("userId")
	r := base.Res{
		Code: base.ResCodeOk,
	}
	if userId == "" {
		r.Code = base.ResCodeRequestParamsError
		r.Message = "用户Id不能为空"
		ctx.JSON(&r)
		return
	}
	service := services.GetAccountService()
	account := service.GetEnvelopeAccountByUserId(userId)
	r.Data = account
	ctx.JSON(&r)
}

//查询账户信息接口：/v1/account/get
func getAccountHandler(ctx iris.Context) {
	accountNo := ctx.URLParam("accountNo")
	r := base.Res{
		Code: base.ResCodeOk,
	}
	if accountNo == "" {
		r.Code=base.ResCodeRequestParamsError
		r.Message= "账户编号不能为空"
		ctx.JSON(&r)
		return
	}
	service := services.GetAccountService()
	account := service.GetAccount(accountNo)
	r.Data = account
	ctx.JSON(&r)
}
