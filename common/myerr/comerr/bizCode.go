package comerr

import (
	"frozen-go-cms/myerr"
)

var (
	// 一般性错误
	TokenInvalid      = myerr.NewBusinessCode(1001, "token invalid", myerr.BusinessData{})
	ExternalIdNoExist = myerr.NewBusinessCode(1003, "externalId no exist", myerr.BusinessData{})
	CodeNoExist       = myerr.NewBusinessCode(1005, "code no exist", myerr.BusinessData{})
	ParaMissing       = myerr.NewBusinessCode(1006, "parameter missing", myerr.BusinessData{})
	InvalidParameter  = myerr.NewBusinessCode(1009, "Invalid parameter", myerr.BusinessData{})
	IncorrectState    = myerr.NewBusinessCode(1013, "Incorrect state", myerr.BusinessData{})
	TransactionFailed = myerr.NewBusinessCode(1014, "Transaction failed", myerr.BusinessData{})
	ReqTooFrequent    = myerr.NewBusinessCode(1018, "Requests are too frequent", myerr.BusinessData{})

	// 钻石
	DiamondNoEnough      = myerr.NewBusinessCode(4000, "Insufficient diamonds", myerr.BusinessData{})
	DiamondFrequency     = myerr.NewBusinessCode(4001, "Diamond operation frequency too high", myerr.BusinessData{})
	DiamondAccountFrozen = myerr.NewBusinessCode(4004, "Diamond Account Frozen", myerr.BusinessData{})

	AlreadyHasCp   = myerr.NewBusinessCode(19001, "The opponent already has CP", myerr.BusinessData{})
	AlreadyExpired = myerr.NewBusinessCode(19002, "Already expired", myerr.BusinessData{})
)
