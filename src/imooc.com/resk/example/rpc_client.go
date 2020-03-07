package main

import (
	"github.com/shopspring/decimal"
	"github.com/sirupsen/logrus"
	"net/rpc"
	"resk-projects/services"
)

func main() {
	c, err := rpc.Dial("tcp", ":8084")
	if err != nil {
		logrus.Panic(err)
	}
	in := services.RedEnvelopeSendingDTO{
		Amount:       decimal.NewFromFloat(1),
		UserId:       "1YiOFIqvGq3vNfZgnx0dpgVH3nK",
		Username:     "测试用户",
		EnvelopeType: services.GeneralEnvelopType,
		Quantity:     2,
		Blessing:     "",
	}

	out := &services.RedEnvelopeActivity{}

	err = c.Call("EnvelopeRpc.SendOut", in, &out)
	if err!= nil {
		logrus.Panic(err)
	}
	logrus.Infof("%+v",out)

}
