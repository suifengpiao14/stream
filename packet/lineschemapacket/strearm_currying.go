package lineschemapacket

import (
	"context"
	"encoding/json"
	"strings"

	"github.com/pkg/errors"
	"github.com/spf13/cast"
	"github.com/suifengpiao14/lineschema"
	"github.com/suifengpiao14/stream"
	"github.com/xeipuuv/gojsonschema"
)

type _MergeDefaultPacketHandler struct {
	BeforeDefaultJson string `json:"beforeDefaultJson"`
	AfterDefaultJson  string `json:"afterDefaultJson"`
}

const PACKETHANDLER_NAME_MergeDefaultPacketHandler = "github.com/suifengpiao14/stream/packet/lineschemapacket/_MergeDefaultPacketHandler"

func NewMergeDefaultHandler(beforeDefaultJson string, afterDefaultJson string) (packet stream.PacketHandlerI) {
	return &_MergeDefaultPacketHandler{
		BeforeDefaultJson: beforeDefaultJson,
		AfterDefaultJson:  afterDefaultJson,
	}
}

func (packet *_MergeDefaultPacketHandler) Name() string {
	return PACKETHANDLER_NAME_MergeDefaultPacketHandler
}

func (packet *_MergeDefaultPacketHandler) Description() string {
	return `合并json默认值`
}

func (packet *_MergeDefaultPacketHandler) String() string {

	return stream.JsonString(packet)
}

func (packet *_MergeDefaultPacketHandler) Before(ctx context.Context, input []byte) (newCtx context.Context, out []byte, err error) {

	return MakeMergeDefaultHandler([]byte(packet.BeforeDefaultJson))(ctx, input)
}

func (packet *_MergeDefaultPacketHandler) After(ctx context.Context, input []byte) (newCtx context.Context, out []byte, err error) {
	return MakeMergeDefaultHandler([]byte(packet.AfterDefaultJson))(ctx, input)
}

func MakeMergeDefaultHandler(defaultJson []byte) (fn stream.HandlerFn) {
	return func(ctx context.Context, input []byte) (newCtx context.Context, out []byte, err error) {
		if len(defaultJson) == 0 {
			return ctx, input, err
		}
		newInput, err := lineschema.MergeDefault(input, defaultJson)
		if err != nil {
			err = errors.WithMessage(err, "merge default value error")
			return ctx, nil, err
		}

		return ctx, newInput, nil
	}
}

type _ValidatePacketHandler struct {
	BeforeJsonSchema     string `json:"beforeJsonschema"`
	AfterJsonSchema      string `json:"afterJsonschema"`
	beforeValidateLoader gojsonschema.JSONLoader
	afterValidateLoader  gojsonschema.JSONLoader
}

const PACKETHANDLER_NAME_ValidatePacket = "github.com/suifengpiao14/stream/packet/lineschemapacket/_ValidatePacketHandler"

func NewValidatePacketHandler(beforeJsonschema string, afterJsonschema string, beforeValidateLoader gojsonschema.JSONLoader, afterValidateLoader gojsonschema.JSONLoader) (packet stream.PacketHandlerI) {
	return &_ValidatePacketHandler{
		BeforeJsonSchema:     beforeJsonschema,
		AfterJsonSchema:      afterJsonschema,
		beforeValidateLoader: beforeValidateLoader,
		afterValidateLoader:  afterValidateLoader,
	}
}

func (packet *_ValidatePacketHandler) Name() string {
	return PACKETHANDLER_NAME_ValidatePacket
}
func (packet *_ValidatePacketHandler) Description() string {
	return ``
}

func (packet *_ValidatePacketHandler) String() string {
	b, _ := json.Marshal(packet)
	s := string(b)
	return s
}
func (packet *_ValidatePacketHandler) Before(ctx context.Context, input []byte) (newCtx context.Context, out []byte, err error) {
	return MakeValidateHandlerFn(packet.beforeValidateLoader)(ctx, input)

}

func (packet *_ValidatePacketHandler) After(ctx context.Context, input []byte) (newCtx context.Context, out []byte, err error) {
	return MakeValidateHandlerFn(packet.afterValidateLoader)(ctx, input)
}

func MakeValidateHandlerFn(validateLoader gojsonschema.JSONLoader) (fn stream.HandlerFn) {
	return func(ctx context.Context, input []byte) (newCtx context.Context, out []byte, err error) {
		if validateLoader == nil {
			return ctx, input, nil
		}
		if len(input) == 0 { // 填充默认格式
			input = []byte("{}")
			jInterface, err := validateLoader.LoadJSON()
			if err != nil {
				return ctx, nil, err
			}
			if m, ok := jInterface.(map[string]any); ok {
				if typ, ok := m["type"]; ok {
					typS := cast.ToString(typ)
					if strings.EqualFold(typS, "array") {
						input = []byte("[]")
					}
				}
			}
		}

		err = lineschema.Validate(input, validateLoader)
		if err != nil {
			return ctx, nil, err
		}
		return ctx, input, nil
	}
}

type _TransferTypeFormatPacketHandler struct {
	BeforePathMap string `json:"beforePathMap"`
	AfterPathMap  string `json:"afterPathMap"`
}

const PACKETHANDLER_NAME_TransferTypeFormatPacket = "github.com/suifengpiao14/stream/packet/lineschemapacket/_TransferTypeFormatPacketHandler"

func NewTransferPacketHandler(beforePathMap string, afterPathMap string) (packet stream.PacketHandlerI) {
	return &_TransferTypeFormatPacketHandler{
		BeforePathMap: beforePathMap,
		AfterPathMap:  afterPathMap,
	}
}

func (packet *_TransferTypeFormatPacketHandler) Name() string {

	return PACKETHANDLER_NAME_TransferTypeFormatPacket
}

func (packet *_TransferTypeFormatPacketHandler) Description() string {

	return `json数据转换,由type类型转成format格式,经过后续处理完又将输出从format格式转为type格式,适用于服务端接口`
}

func (packet *_TransferTypeFormatPacketHandler) String() string {

	return stream.JsonString(packet)
}

func (packet *_TransferTypeFormatPacketHandler) Before(ctx context.Context, input []byte) (newCtx context.Context, out []byte, err error) {
	return MakeTransferHandler(packet.BeforePathMap)(ctx, input)
}

func (packet *_TransferTypeFormatPacketHandler) After(ctx context.Context, input []byte) (newCtx context.Context, out []byte, err error) {
	return MakeTransferHandler(packet.AfterPathMap)(ctx, input)
}

func MakeTransferHandler(pathMap string) (fn stream.HandlerFn) {
	return func(ctx context.Context, input []byte) (newCtx context.Context, out []byte, err error) {
		out = lineschema.ConvertFomat(input, pathMap)
		return ctx, out, nil
	}
}
