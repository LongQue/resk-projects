package envelopes

import (
	"github.com/sirupsen/logrus"
	"github.com/tietang/dbx"
)

type RedEnvelopeGoodsDao struct {
	runner *dbx.TxRunner
}

//插入
func (dao *RedEnvelopeGoodsDao)Insert(po *RedEnvelopeGoods)  (int64,error){
	rs, err := dao.runner.Insert(po)
	if err!=nil{
		return 0,err
	}
	return rs.LastInsertId()
}

//查询，根据红包编号
func (dao*RedEnvelopeGoodsDao)GetOne(envelopeNo string)  *RedEnvelopeGoods{
	po := &RedEnvelopeGoods{
		EnvelopeNo: envelopeNo,
	}
	ok,err := dao.runner.GetOne(po)
	if err!=nil||!ok{
		logrus.Error(err)
		return nil
	}
	return po
}
//更新红包余额和数量

//更新订单状态

//过期，把过期的所有红包都查询出来，分页，limit offset  size