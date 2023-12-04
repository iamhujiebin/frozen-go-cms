//go:build !windows
// +build !windows

package main

import (
	"frozen-go-cms/wasm/tools"
	"syscall/js"
)

func main() {
	done := make(chan struct{})
	global := js.Global()
	global.Set("wasmHumanReadableTimediff", js.FuncOf(humanReadableTimediff))
	global.Set("wasmUnixTimeConverter", js.FuncOf(unixTimeConverter))
	global.Set("wasmDateTimeConverter", js.FuncOf(dateTimeConverter))
	global.Set("wasmEncodeDecode", js.FuncOf(encodeDecode))
	global.Set("wasmGenerateQRCode", js.FuncOf(generateQRCode))
	<-done
}

func encodeDecode(this js.Value, args []js.Value) interface{} {
	if len(args) != 3 {
		return "ERROR: number of arguments doesn't match"
	}
	value, action, _type := args[0].String(), args[1].String(), args[2].String()
	result := tools.EncodeDecode(value, action, _type)
	return result
}

func generateQRCode(this js.Value, args []js.Value) interface{} {
	if len(args) != 1 {
		return "ERROR: number of arguments doesn't match"
	}
	result := tools.GenerateQRCode(args[0].String())
	return result
}

func unixTimeConverter(this js.Value, args []js.Value) interface{} {
	if len(args) != 1 {
		return "ERROR: number of arguments doesn't match"
	}
	unixTime := int64(args[0].Int())
	return tools.UnixTimeConverter(unixTime)
}

func dateTimeConverter(this js.Value, args []js.Value) interface{} {
	if len(args) != 1 {
		return "ERROR: number of arguments doesn't match"
	}
	dateTime := args[0].String()
	return tools.DateTimeConverter(dateTime)
}

func humanReadableTimediff(this js.Value, args []js.Value) interface{} {
	if len(args) != 1 {
		return "ERROR: number of arguments doesn't match"
	}
	timediff := float32(args[0].Float())
	return tools.HumanReadableTimeDiff(timediff)
}
