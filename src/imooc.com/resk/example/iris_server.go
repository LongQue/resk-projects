package main

import (
	"fmt"
	"github.com/kataras/iris"
	"github.com/sirupsen/logrus"
	"net/http"
	"strconv"
)

func main() {
	app := iris.Default()
	app.Get("/hello", func(ctx iris.Context) {
		ctx.WriteString("hello,world!")
	})
	v1 := app.Party("/v1")
	v1.Use(func(ctx iris.Context) {
		logrus.Info("自定义中间件")
		ctx.Next()
	})
	v1.Get("/users/{id:uint64 min(2)}", func(ctx iris.Context) {
		id := ctx.Params().GetUint64Default("id", 0)
		ctx.WriteString(strconv.Itoa(int(id)))
	})

	v1.Get("/orders/{action:string prefix(a_)}", func(ctx iris.Context) {
		a := ctx.Params().Get("action")
		ctx.WriteString(a)
	})
	//覆盖所有状态码错误，谁写在后面谁优先级高
	app.OnAnyErrorCode(func(ctx iris.Context) {
		ctx.WriteString("错误")
	})
	//根据状态码返回错误信息
	app.OnErrorCode(http.StatusInternalServerError,
		func(ctx iris.Context) {
			ctx.WriteString("服务器内部错误")
		})
	err := app.Run(iris.Addr(":8082"))
	fmt.Println(err)
}
