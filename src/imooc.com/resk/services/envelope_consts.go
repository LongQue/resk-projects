package services

const (
	DefaultBlessing = "恭喜发财"
)

//订单类型：发布单、退款单
type OrderType int

const (
	OrderTypeSending OrderType = 1
	OrderTypeRefund  OrderType = 2
)

//支付:未支付，支付中，已支付，支付失败
//退款:未退款，退款中，已退款，退款失败
type PayStatus int

const (
	PayNothing PayStatus = 1
	Paying     PayStatus = 2
	Payed      PayStatus = 3
	PayFailure PayStatus = 4

	RefundNothing PayStatus = 61
	Refunding     PayStatus = 62
	Refunded      PayStatus = 63
	RefundFailure PayStatus = 64
)

//红包订单状态：创建、发布、过期、失效
type OrderStatus int

const (
	OrderCreate                  OrderStatus = 1
	OrderSending                 OrderStatus = 2
	OrderExpired                 OrderStatus = 3
	OrderDisabled                OrderStatus = 4
	OrderExpiredRefundSuccessful OrderStatus = 5
	OrderExpiredRefundFailure    OrderStatus = 6
)

//红包类型：普通红包、运气红包
type EnvelopeType int

const (
	GeneralEnvelopType EnvelopeType = 1
	LuckyEnvelopeType  EnvelopeType = 2
)
