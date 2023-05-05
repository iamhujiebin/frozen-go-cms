package bizerr

import (
	"frozen-go-cms/myerr"
)

var (
	// 一般性错误
	AuthFail     = myerr.NewBusinessCode(1000, "auth fail", myerr.BusinessData{})
	TokenInvalid = myerr.NewBusinessCode(1001, "token invalid", myerr.BusinessData{})
	ParaMissing  = myerr.NewBusinessCode(1006, "parameter missing", myerr.BusinessData{})
)
