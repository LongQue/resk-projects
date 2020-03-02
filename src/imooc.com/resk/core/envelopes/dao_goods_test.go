package envelopes

import (
	"database/sql"
	"github.com/segmentio/ksuid"
	"github.com/shopspring/decimal"
	. "github.com/smartystreets/goconvey/convey"
	"github.com/tietang/dbx"
	"resk-projects/infra/base"
	"resk-projects/services"
	"testing"
	"time"
)

//1、红包商品数据写入
func TestRedEnvelopeGoodsDao_Insert(t *testing.T) {
	err := base.Tx(func(runner *dbx.TxRunner) error {
		dao := &RedEnvelopeGoodsDao{
			runner: runner,
		}
		Convey("红包商品数据写入", t, func() {
			a := &RedEnvelopeGoods{
				EnvelopeNo:     ksuid.New().Next().String(),
				EnvelopeType:   services.GeneralEnvelopType,
				Username:       sql.NullString{String: "测试用户", Valid: true},
				UserId:         ksuid.New().Next().String(),
				Blessing:       sql.NullString{String: "祝福语", Valid: true},
				Amount:         decimal.NewFromFloat(60),
				AmountOne:      decimal.NewFromFloat(6),
				Quantity:       10,
				RemainAmount:   decimal.NewFromFloat(36),
				RemainQuantity: 6,
				ExpiredAt:      time.Now(),
				Status:         services.OrderSending,
				OrderType:      services.OrderTypeSending,
				PayStatus:      services.Paying,
			}
		})

	})
}

//2、更新红包剩余金额和数量
