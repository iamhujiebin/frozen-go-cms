package tools

import (
	"encoding/base64"
	"fmt"
	"net/url"
	"time"

	"github.com/skip2/go-qrcode"
)

type requestStruct struct {
	Value, Action, Type string
}

func EncodeDecode(value, action, _type string) string {
	t := requestStruct{value, action, _type}
	var result string
	switch t {
	case requestStruct{t.Value, "encode", "url"}:
		result = url.PathEscape(t.Value)
	case requestStruct{t.Value, "decode", "url"}:
		result, _ = url.PathUnescape(t.Value)
	case requestStruct{t.Value, "encode", "base64"}:
		result = base64.StdEncoding.EncodeToString([]byte(t.Value))
	case requestStruct{t.Value, "decode", "base64"}:
		decoded, err := base64.StdEncoding.DecodeString(t.Value)
		if err != nil {
			result = "ERROR: " + err.Error()
		} else {
			result = string(decoded)
		}
	default:
		result = "ERROR: 'Action' should 'encode' or 'decode' and 'Type' should 'url' or 'base64'"
	}
	return result
}

func GenerateQRCode(val string) string {
	var png []byte
	png, err := qrcode.Encode(val, qrcode.Medium, 256)
	if err != nil {
		return "ERROR: error during generating QR Code"
	}
	pngBase64 := base64.StdEncoding.EncodeToString(png)
	return pngBase64
}

func UnixTimeConverter(unixTime int64) string {
	value := time.Unix(unixTime, 0).Local()
	return fmt.Sprintf("%s", value)
}

const minLength float32 = 60
const hourLength = minLength * 60
const dayLength = hourLength * 24
const monthLength = dayLength * 31
const yearLength = dayLength * 365

func HumanReadableTimeDiff(timediff float32) interface{} {

	var agoOrLater string
	if timediff >= 0 {
		agoOrLater = "ago"
	} else {
		timediff = -timediff
		agoOrLater = "later"
	}

	var humanReadableResult string

	if timediff < minLength {
		humanReadableResult = fmt.Sprintf("%.1f seconds %s", timediff, agoOrLater)
	} else if timediff < hourLength {
		humanReadableResult = fmt.Sprintf("%.1f minutes %s", timediff/minLength, agoOrLater)
	} else if timediff < dayLength {
		humanReadableResult = fmt.Sprintf("%.1f hours %s", timediff/hourLength, agoOrLater)
	} else if timediff < monthLength {
		humanReadableResult = fmt.Sprintf("%.1f days %s", timediff/dayLength, agoOrLater)
	} else if timediff < yearLength {
		humanReadableResult = fmt.Sprintf("%.1f months %s", timediff/monthLength, agoOrLater)
	} else {
		humanReadableResult = fmt.Sprintf("%.1f years %s", timediff/yearLength, agoOrLater)
	}

	return humanReadableResult
}
