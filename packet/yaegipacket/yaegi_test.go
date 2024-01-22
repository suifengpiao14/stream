package yaegipacket_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/suifengpiao14/stream/packet/yaegipacket"
)

func TestNewCurlHookYaegi(t *testing.T) {
	t.Run("has BeforeFn and AfterFn", func(t *testing.T) {
		dynamic := `
		package curlhook

import (
	"github.com/tidwall/gjson"
	"github.com/tidwall/sjson"
)

func BeforeFn(input string) (output string, err error) {
	timestamps := gjson.Get(input, "body._head._timestamps").String()
	_ = timestamps
	input, err = sjson.Set(input, "body._head._timestamps", "1111111111111111")
	if err != nil {
		return "", err
	}
	return input, nil
}
func AfterFn(input string) (output string, err error) {
	return input, nil
}

	
		`
		_, err := yaegipacket.NewCurlHookYaegi(dynamic)
		require.NoError(t, err)
	})

	t.Run("only BeforFn", func(t *testing.T) {
		dynamic := `
		package curlhook

import (
	"github.com/tidwall/gjson"
	"github.com/tidwall/sjson"
)

func BeforeFn(input string) (output string, err error) {
	timestamps := gjson.Get(input, "body._head._timestamps").String()
	_ = timestamps
	input, err = sjson.Set(input, "body._head._timestamps", "1111111111111111")
	if err != nil {
		return "", err
	}
	return input, nil
}


		`
		_, err := yaegipacket.NewCurlHookYaegi(dynamic)
		require.NoError(t, err)
	})

	t.Run("only AfterFn", func(t *testing.T) {
		dynamic := `
		package curlhook
		func AfterFn(input string) (output string, err error) {
			return input, nil
		}
		
		`
		_, err := yaegipacket.NewCurlHookYaegi(dynamic)
		require.NoError(t, err)
	})

}
