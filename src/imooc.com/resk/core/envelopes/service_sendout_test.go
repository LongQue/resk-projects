package envelopes

import (
	"github.com/segmentio/ksuid"
	"github.com/shopspring/decimal"
	. "github.com/smartystreets/goconvey/convey"
	"resk-projects/services"
	"testing"
)

func TestRedEnvelopeService_SendOut(t *testing.T) {
	//发红包人的红包资金账户
	ac := services.GetAccountService()
	account := services.AccountCreatedDTO{
		UserId:       ksuid.New().Next().String(),
		Username:     "测试用户",
		AccountName:  "测试账户",
		AccountType:  int(services.EnvelopeAccountType),
		CurrencyCode: "CNY",
		Amount:       "100",
	}
	re := services.GetRedEnvelopeService()
	Convey("准备资金账户",t, func() {
		acDTO, err := ac.CreateAccount(account)
		So(err, ShouldBeNil)
		So(acDTO,ShouldNotBeNil)
	})
	Convey("发送红包",t, func() {
		goods:= services.RedEnvelopeSendingDTO{
			EnvelopeType: services.GeneralEnvelopType,
			Username:     account.Username,
			UserId:       account.UserId,
			Blessing:     services.DefaultBlessing,
			Amount:       decimal.NewFromFloat(8.88),
			Quantity:     10,
		}

		Convey("发普通红包", func() {
			at, err := re.SendOut(goods)
			So(err,ShouldBeNil)
			So(at, ShouldNotBeNil)
			So(at.Link, ShouldNotBeEmpty)
			So(at.RedEnvelopeGoodsDTO,ShouldNotBeNil)
			//验证每一个属性
			dto:=at.RedEnvelopeGoodsDTO
			So(dto.Username,ShouldEqual,goods.Username)
			So(dto.UserId,ShouldEqual,goods.UserId)
			So(dto.Quantity,ShouldEqual,goods.Quantity)
			q := decimal.NewFromFloat(float64(dto.Quantity))
			So(dto.Amount.String(),ShouldEqual,goods.Amount.Mul(q).String())
			So(dto.EnvelopeType,ShouldEqual,goods.EnvelopeType)
		})
		goods.EnvelopeType=services.LuckyEnvelopeType
		goods.Amount=decimal.NewFromFloat(88.8)
		Convey("发运气红包", func() {
			at, err := re.SendOut(goods)
			So(err,ShouldBeNil)
			So(at, ShouldNotBeNil)
			So(at.Link, ShouldNotBeEmpty)
			So(at.RedEnvelopeGoodsDTO,ShouldNotBeNil)
			//验证每一个属性
			dto:=at.RedEnvelopeGoodsDTO
			So(dto.Username,ShouldEqual,goods.Username)
			So(dto.UserId,ShouldEqual,goods.UserId)
			So(dto.Quantity,ShouldEqual,goods.Quantity)
			So(dto.Amount.String(),ShouldEqual,goods.Amount.String())
			So(dto.EnvelopeType,ShouldEqual,goods.EnvelopeType)
		})
	})
}