package bizerr

import (
	"frozen-go-cms/myerr"
)

var (
	// 一般性错误
	AuthFail    = myerr.NewBusinessCode(1000, "auth fail", myerr.BusinessData{})
	ParaMissing = myerr.NewBusinessCode(1006, "parameter missing", myerr.BusinessData{})
)
