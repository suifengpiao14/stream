package packet

import (
	"context"
	"encoding/json"

	"github.com/suifengpiao14/packethandler"
)

type _JsonMarshalUnMarshalPacket struct {
	dataProvider interface{}
	dataReceiver interface{}
}

const PACKETHANDLER_NAME_JsonMarshalUnMarshalPacket = "github.com/suifengpiao14/stream/packet/_JsonMarshalUnMarshalPacket"

// NewJsonMarshalUnMarshalPacket 结构体转字节再转结构体
func NewJsonMarshalUnMarshalPacket(dataProvider interface{}, dataReceiver interface{}) (pack packethandler.PacketHandlerI) {
	return &_JsonMarshalUnMarshalPacket{
		dataProvider: dataProvider,
		dataReceiver: dataReceiver,
	}
}

func (pack *_JsonMarshalUnMarshalPacket) Name() string {
	return PACKETHANDLER_NAME_JsonMarshalUnMarshalPacket
}
func (pack *_JsonMarshalUnMarshalPacket) Description() string {
	return "struct -> []byte -> struct"
}

func (pack *_JsonMarshalUnMarshalPacket) String() string {
	return ""
}

func (pack *_JsonMarshalUnMarshalPacket) Before(ctx context.Context, input []byte) (newCtx context.Context, out []byte, err error) {
	if pack.dataProvider == nil {
		return ctx, input, nil
	}
	out, err = json.Marshal(pack.dataProvider)
	return ctx, out, err
}

func (pack *_JsonMarshalUnMarshalPacket) After(ctx context.Context, input []byte) (newCtx context.Context, out []byte, err error) {
	if pack.dataReceiver == nil {
		return ctx, input, nil
	}
	if len(input) == 0 {
		return ctx, input, nil
	}
	err = json.Unmarshal(input, pack.dataReceiver)
	return ctx, nil, err
}

type _JsonUnmarshalMarshalPacket struct {
	_JsonMarshalUnMarshalPacket
}

const PACKETHANDLER_NAME_JsonUnmarshalMarshalPacket = "github.com/suifengpiao14/stream/packet/_JsonUnmarshalMarshalPacket"

// NewJsonUnmarshalMarshalPacket 字节转结构体再转字节
func NewJsonUnmarshalMarshalPacket(dataReceiver interface{}, dataProvider interface{}) (pack packethandler.PacketHandlerI) {
	return &_JsonUnmarshalMarshalPacket{
		_JsonMarshalUnMarshalPacket: _JsonMarshalUnMarshalPacket{
			dataProvider: dataProvider,
			dataReceiver: dataReceiver,
		},
	}
}

func (pack *_JsonUnmarshalMarshalPacket) Name() string {
	return PACKETHANDLER_NAME_JsonUnmarshalMarshalPacket
}
func (pack *_JsonUnmarshalMarshalPacket) Description() string {
	return "[]byte -> struct -> []byte"
}

func (pack *_JsonUnmarshalMarshalPacket) Before(ctx context.Context, input []byte) (newCtx context.Context, out []byte, err error) {
	return pack._JsonMarshalUnMarshalPacket.After(ctx, input)
}

func (pack *_JsonUnmarshalMarshalPacket) After(ctx context.Context, input []byte) (newCtx context.Context, out []byte, err error) {
	return pack._JsonMarshalUnMarshalPacket.Before(ctx, input)
}
