package envelopes

import (
	"context"
	"errors"
	"github.com/shopspring/decimal"
	log "github.com/sirupsen/logrus"
	"github.com/tietang/dbx"
	"resk-projects/infra/base"
	"resk-projects/services"
)

const (
	pageSize = 100
)

type ExpiredEnvelopeDomain struct {
	expiredGoods []RedEnvelopeGoods
	offset       int
}

//查询出过期红包
func (e *ExpiredEnvelopeDomain) Next() (ok bool) {
	base.Tx(func(runner *dbx.TxRunner) error {
		dao := RedEnvelopeGoodsDao{runner: runner}
		e.expiredGoods = dao.FindExpired(e.offset, pageSize)
		if len(e.expiredGoods) > 0 {
			e.offset += len(e.expiredGoods)
			ok = true
		}
		return nil
	})
	return ok
}

func (e *ExpiredEnvelopeDomain) Expired() (err error) {
	for e.Next() {
		for _, g := range e.expiredGoods {
			log.Debug("过期红包退款开始:,%+v", g)
			err := e.ExpiredOne(g)
			if err != nil {
				log.Error(err)
			}
			log.Debug("过期红包退款结束:,%+v", g)
		}
	}
	return nil
}

//发起退款流程
func (e *ExpiredEnvelopeDomain) ExpiredOne(goods RedEnvelopeGoods) (err error) {
	//创建一个退款订单
	refund := goods
	refund.OrderType = services.OrderTypeRefund
	refund.RemainAmount = goods.RemainAmount.Mul(decimal.NewFromFloat(-1))
	refund.RemainQuantity = -goods.RemainQuantity
	refund.Status = services.OrderExpired
	refund.PayStatus = services.Refunding
	refund.OriginEnvelopeNo = goods.EnvelopeNo
	refund.EnvelopeNo = ""
	domain := goodsDomain{RedEnvelopeGoods: refund}
	domain.createEnvelopeNo()
	//调用资金账户接口进行转账
	err = base.Tx(func(runner *dbx.TxRunner) error {
		txCtx := base.WithValueContext(context.Background(), runner)
		id, err := domain.Save(txCtx)
		if err != nil || id == 0 {
			return errors.New("创建退款订单失败")
		}
		//修改原订单状态
		dao := RedEnvelopeGoodsDao{runner: runner}
		rows, err := dao.UpdateOrderStatus(goods.EnvelopeNo, services.OrderExpired)
		if err != nil || rows == 0 {
			return errors.New("更新原订单状态失败")
		}
		return nil
	})
	if err != nil {
		return
	}
	//调用资金账户接口进行转账
	systemAccount := base.GetSystemAccount()
	account := services.GetAccountService().GetEnvelopeAccountByUserId(goods.UserId)
	if account == nil {
		return errors.New("没有找到该用户的红包资金账户:" + goods.UserId)
	}
	body := services.TradeParticipator{
		AccountNo: systemAccount.AccountNo,
		UserId:    systemAccount.UserId,
		Username:  systemAccount.Username,
	}
	target := services.TradeParticipator{
		AccountNo: account.AccountNo,
		UserId:    account.UserId,
		Username:  account.Username,
	}
	transfer := services.AccountTransferDTO{
		TradeNo:     refund.EnvelopeNo,
		TradeBody:   body,
		TradeTarget: target,
		Amount:      goods.RemainAmount,
		ChangeType:  services.EnvelopExpiredRefund,
		ChangeFlag:  services.FlagTransferIn,
		Decs:        "红包过期退款：" + goods.EnvelopeNo,
	}
	status, err := services.GetAccountService().Transfer(transfer)
	if status != services.TransferStatusSuccess {
		return err
	}
	//修改原订单状态
	err = base.Tx(func(runner *dbx.TxRunner) error {
		dao := RedEnvelopeGoodsDao{runner: runner}
		//修改原订单状态
		rows, err := dao.UpdateOrderStatus(goods.EnvelopeNo, services.OrderExpiredRefundSuccessful)
		if err != nil || rows == 0 {
			return errors.New("更新原订单状态失败")
		}
		//修改退款订单状态
		rows, err = dao.UpdateOrderStatus(goods.EnvelopeNo, services.OrderExpiredRefundSuccessful)
		if err != nil || rows == 0 {
			return errors.New("更新退款订单状态失败")
		}
		return nil
	})
	if err != nil {
		log.Error(err)
		return err
	}
	return nil

}
