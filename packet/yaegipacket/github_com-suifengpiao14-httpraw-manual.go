package yaegipacket

import (
	"reflect"

	"github.com/suifengpiao14/httpraw"
)

func init() {
	Symbols["github.com/suifengpiao14/httpraw/httpraw"] = map[string]reflect.Value{
		// function, constant and variable definitions
		"CurrentTime":     reflect.ValueOf(httpraw.CurrentTime),
		"Fen2yuan":        reflect.ValueOf(httpraw.Fen2yuan),
		"GetMD5LOWER":     reflect.ValueOf(httpraw.GetMD5LOWER),
		"JsonCompact":     reflect.ValueOf(httpraw.JsonCompact),
		"TimestampSecond": reflect.ValueOf(httpraw.TimestampSecond),
		"WithDefault":     reflect.ValueOf(httpraw.WithDefault),
		"WithEmptyStr":    reflect.ValueOf(httpraw.WithEmptyStr),
		"WithZeroNumber":  reflect.ValueOf(httpraw.WithZeroNumber),
		"Xid":             reflect.ValueOf(httpraw.Xid),
		"ZeroTime":        reflect.ValueOf(httpraw.ZeroTime),
		"RequestDTO":      reflect.ValueOf((*httpraw.RequestDTO)(nil)),
	}
}
