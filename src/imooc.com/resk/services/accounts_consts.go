package services

//转账状态
type TransferStatus int8

const (
	//转账失败
	TransferStatusFailure TransferStatus = -1
	//余额不足
	TransferStatusSufficientFunds TransferStatus = -0
	//转账成功
	TransferStatusSuccess TransferStatus = 1
)

//转账类型 0：创建账户  >=1 进账 <=-1 支出
type ChangeType int8

const (
	//创建账户
	AccountCreated ChangeType = 0
	//储值
	AccountStoreValue ChangeType = 1
	//红包资金支出
	EnvelopeOutgoing ChangeType = -2
	//红包资金收入
	EnvelopeIncoming ChangeType = 2
	//红包过期退款
	EnvelopExpiredRefund ChangeType = 3
)

type ChangeFlag int8

const (
	//创建账户=0
	FlagAccountCreate ChangeFlag = 0
	//支出=-1
	FlagTransferOut ChangeFlag = -1
	//收入=1
	FlagTransferIn ChangeFlag = 1
)

//账户类型
type AccountType int8

const (
	EnvelopeAccountType           AccountType = 1
	SystemEnvelopeTypeAccountType AccountType = 2
)
