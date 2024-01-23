package yaegipacket

import (
	"context"
	"fmt"
	"reflect"
	"strings"

	"github.com/pkg/errors"
	_ "github.com/spf13/cast"
	"github.com/suifengpiao14/stream"
	_ "github.com/tidwall/gjson"
	_ "github.com/tidwall/sjson"
	"github.com/traefik/yaegi/interp"
	"github.com/traefik/yaegi/stdlib"
)

var Symbols = stdlib.Symbols
var CURLHookBeforeFnPoint = "curlhook.BeforeFn"
var CURLHookAfterFnPoint = "curlhook.AfterFn"

// 动态脚本构造函数格式,兼容js,tengo,php脚本定义，所以 定义成最通用的string
type HookFn func(input string) (output string, err error)

const (
	undefined_selector_error_prefix = "undefined selector: "
)

// Validate  用于往数据库预先写入动态脚本时验证合法性
func Validate(dynamicScript string) (err error) {
	_, err = NewCurlHookYaegi(dynamicScript)
	return err
}

type YaegiHook struct {
	dynamicScript string
	DynamicBefore stream.HandlerFn
	DynamicAfter  stream.HandlerFn
}

func NewCurlHookYaegi(dynamicScript string) (yaegiHook *YaegiHook, err error) {
	var (
		beforeFn       HookFn
		afterFn        HookFn
		beforeFnExists = true
		afterFnExists  = true
	)
	yaegiHook = &YaegiHook{
		dynamicScript: dynamicScript,
	}

	// 解析动态脚本
	interpreter := interp.New(interp.Options{})
	interpreter.Use(stdlib.Symbols)

	interpreter.Use(Symbols) //注册当前包结构体

	_, err = interpreter.Eval(dynamicScript)
	if err != nil {
		err = errors.WithMessage(err, "init dynamic go script error")
		return nil, err
	}
	fnT := reflect.TypeOf((HookFn)(nil))
	beforeFnV, beforeFnExists, err := getFn(interpreter, CURLHookBeforeFnPoint, fnT)
	if err != nil {
		return nil, err
	}

	if beforeFnExists {
		beforeFn = beforeFnV.Interface().(HookFn)
	}

	afterFnV, afterFnExists, err := getFn(interpreter, CURLHookAfterFnPoint, fnT)
	if err != nil {
		return nil, err
	}
	if afterFnExists {
		afterFn = afterFnV.Interface().(HookFn)
	}
	yaegiHook.DynamicBefore = func(ctx context.Context, input []byte) (newCtx context.Context, out []byte, err error) {
		if beforeFn == nil {
			return stream.EmptyHandlerFn(ctx, input)
		}
		inputS := string(input)
		outS, err := beforeFn(inputS)
		out = []byte(outS)
		return ctx, out, err
	}
	yaegiHook.DynamicAfter = func(ctx context.Context, input []byte) (newCtx context.Context, out []byte, err error) {
		if afterFn == nil {
			return stream.EmptyHandlerFn(ctx, input)
		}
		inputS := string(input)
		outS, err := afterFn(inputS)
		out = []byte(outS)
		return ctx, out, err
	}
	return yaegiHook, nil
}

// getFn 从动态脚本中获取特定函数
func getFn(interpreter *interp.Interpreter, selector string, dstType reflect.Type) (fn reflect.Value, exists bool, err error) {
	fnV, err := interpreter.Eval(selector)
	if err != nil && strings.Contains(err.Error(), undefined_selector_error_prefix) { // 不存在当前元素 时 忽略错误，程序容许只动态实现一部分
		err = nil
		return fn, false, nil
	}

	if err != nil {
		err = errors.WithMessage(err, selector)
		return fn, false, err
	}
	if !fnV.CanConvert(dstType) {
		err = errors.Errorf("dynamic func %s ,must can convert to %s", selector, fmt.Sprintf("%s.%s", dstType.PkgPath(), dstType.Name()))
		return fn, true, err
	}
	fn = fnV.Convert(dstType)
	return fn, true, nil
}

//go:generate go install github.com/traefik/yaegi/cmd/yaegi
//go:generate yaegi extract github.com/suifengpiao14/httpraw
//go:generate yaegi extract github.com/tidwall/gjson
//go:generate yaegi extract github.com/tidwall/sjson
//go:generate yaegi extract github.com/spf13/cast
//go:generate yaegi extract github.com/suifengpiao14/stream/packet/yaegipacket/customfunc
//go:generate yaegi extract github.com/syyongx/php2go
//github.com/suifengpiao14/gjsonmodifier
