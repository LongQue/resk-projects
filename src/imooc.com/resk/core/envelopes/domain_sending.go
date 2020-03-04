package envelopes

import (
	"context"
	"github.com/tietang/dbx"
	"path"
	"resk-projects/core/accounts"
	"resk-projects/infra/base"
	"resk-projects/services"
)

//发红包业务领域代码
func (d *goodsDomain) SendOut(goods services.RedEnvelopeGoodsDTO) (activity *services.RedEnvelopeActivity, err error) {
	//创建红包商品
	d.Create(goods)
	//创建活动
	activity = new(services.RedEnvelopeActivity)
	//红包链接 //http:/域名/v1/envelope/link/{id}/
	link := base.GetEnvelopeActivityLink()
	domain := base.GetEnvelopeDomain()
	activity.Link = path.Join(domain, link, d.EnvelopeNo)
	accountDomain := accounts.NewAccountDomain()

	err = base.Tx(func(runner *dbx.TxRunner) error {
		ctx := base.WithValueContext(context.Background(), runner)
		//事务逻辑问题：
		//保存红包商品和红包金额的支付必须保证全部成功或者全部失败

		//保存红包商品
		id, err := d.Save(ctx)
		if id < 0 || err != nil {
			return err
		}
		//红包金额支付
		//1.需要红包中间商的红包资金账户，定义在配置文件中，事前初始化到资金账户表中
		//2。从红包发送人的资金账户中红扣减红包金额
		//3.将扣减的红包总金额转入红包中间商的红包资金账户，资金从红包发送人的账户扣除
		body := services.TradeParticipator{
			AccountNo: goods.AccountNo,
			UserId:    goods.UserId,
			Username:  goods.Username,
		}
		systemAccount := base.GetSystemAccount()
		target := services.TradeParticipator{
			AccountNo: systemAccount.AccountNo,
			UserId:    systemAccount.UserId,
			Username:  systemAccount.Username,
		}

		transfer := services.AccountTransferDTO{
			TradeBody:   body,
			TradeTarget: target,
			TradeNo:     d.EnvelopeNo,
			Amount:      d.Amount,
			ChangeType:  services.EnvelopeOutgoing,
			ChangeFlag:  services.FlagTransferOut,
			Decs:        "红包金额支付",
		}
		status, err := accountDomain.TransferWithContextTx(ctx, transfer)
		if status == services.TransferStatusSuccess {
			return nil
		}
		//3.将扣减的红包总金额转入红包中间商的红包资金账户
		transfer = services.AccountTransferDTO{
			TradeBody:   target,
			TradeTarget: body,
			TradeNo:     d.EnvelopeNo,
			Amount:      d.Amount,
			ChangeType:  services.EnvelopeOutgoing,
			ChangeFlag:  services.FlagTransferIn,
			Decs:        "红包金额转入",
		}
		status, err = accountDomain.TransferWithContextTx(ctx, transfer)
		if status == services.TransferStatusSuccess {
			return nil
		}
		return err
	})
	if err != nil {
		return nil, err
	}

	//扣减金额没有问题，返回活动
	activity.RedEnvelopeGoodsDTO = *d.RedEnvelopeGoods.ToDTO()

	return activity, err

}
