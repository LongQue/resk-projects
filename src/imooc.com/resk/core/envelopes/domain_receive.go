package envelopes

import (
	"context"
	"database/sql"
	"errors"
	"github.com/shopspring/decimal"
	"github.com/tietang/dbx"
	"resk-projects/core/accounts"
	"resk-projects/infra/algo"
	"resk-projects/infra/base"
	"resk-projects/services"
)

var multiple = decimal.NewFromFloat(100.0)

func (d *goodsDomain) Receive(ctx context.Context, dto services.RedEnvelopeReceiveDTO) (item *services.RedEnvelopeItemDTO, err error) {
	//1.创建红包的订单明细 perCreateItem
	d.preCreateItem(dto)
	//2.查询当前红包的剩余数量和剩余金额信息
	goods := d.Get(dto.EnvelopeNo)
	//3.校验剩余红包和剩余金额，如果没有返回无可用红包
	if goods.RemainQuantity <= 0 || goods.RemainAmount.Cmp(decimal.NewFromFloat(0)) <= 0 {
		return nil, errors.New("没有足够的红包和金额")

	}
	//4.使用红包算法计算红包金额
	nextAmount := d.nextAmount(goods)
	err = base.Tx(func(runner *dbx.TxRunner) error {
		//5.使用乐观锁更新语句，尝试更新剩余数量和剩余金额
		dao := RedEnvelopeGoodsDao{runner: runner}
		rows, err := dao.UpdateBalance(goods.EnvelopeNo, nextAmount)
		// - 失败，返回0，表示无可用红包数量和金额，抢红包失败
		if rows <= 0 || err != nil {
			return errors.New("没有足够的红包和金额了")
		}
		//6.保存订单明细数据
		d.item.Quantity = 1
		d.item.PayStatus = int(services.Paying)
		d.item.AccountNo = dto.AccountNo
		d.item.RemainAmount = goods.RemainAmount.Sub(nextAmount)
		d.item.Amount = nextAmount
		txCtx := base.WithValueContext(ctx, runner)
		_, err = d.item.Save(txCtx)
		if err != nil {
			return err
		}
		//7.将抢到的红包金额从系统红包中间账户转入当前用户的资金账户
		//transfer
		status, err := d.transfer(txCtx, dto)
		if status==services.TransferStatusSuccess {
			return nil
		}
		return err
	})
	return d.item.ToDTO(),err
}

func (d *goodsDomain) transfer(ctx context.Context, dto services.RedEnvelopeReceiveDTO) (status services.TransferStatus, err error) {
	systemAccount := base.GetSystemAccount()
	body := services.TradeParticipator{
		AccountNo: systemAccount.AccountNo,
		UserId:    systemAccount.UserId,
		Username:  systemAccount.Username,
	}
	target := services.TradeParticipator{
		AccountNo: dto.AccountNo,
		UserId:    dto.RecvUserId,
		Username:  dto.RecvUsername,
	}
	transfer := services.AccountTransferDTO{
		TradeNo:     dto.EnvelopeNo,
		TradeBody:   body,
		TradeTarget: target,
		Amount:      d.item.Amount,
		ChangeType:  services.EnvelopeIncoming,
		ChangeFlag:  services.FlagTransferIn,
		Decs:        "红包输入",
	}
	adomain := accounts.NewAccountDomain()
	return adomain.TransferWithContextTx(ctx, transfer)
}

//预创建收红包订单明细
func (d *goodsDomain) preCreateItem(dto services.RedEnvelopeReceiveDTO) {
	d.item.AccountNo = dto.AccountNo
	d.item.EnvelopeNo = dto.EnvelopeNo
	d.item.RecvUsername = sql.NullString{String: dto.RecvUsername}
	d.item.RecvUserId = dto.RecvUserId
	d.item.createItemNo()
}

//计算红包金额
func (d *goodsDomain) nextAmount(goods *RedEnvelopeGoods) (amount decimal.Decimal) {
	if goods.RemainQuantity == 1 {
		return goods.RemainAmount
	}
	if goods.EnvelopeType == services.GeneralEnvelopType {
		return goods.Amount
	}
	if goods.EnvelopeType == services.LuckyEnvelopeType {
		cent := goods.RemainAmount.Mul(multiple).IntPart()
		next := algo.DoubleAverage(int64(d.RemainQuantity), cent)
		amount = decimal.NewFromFloat(float64(next)).Div(multiple)
	}
	return amount

}
