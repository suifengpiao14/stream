package packet

import (
	"context"
	"encoding/json"

	"github.com/suifengpiao14/stream"
)

type _JsonMarshalUnMarshalPacket struct {
	dataProvider interface{}
	dataReceiver interface{}
}

//NewJsonMarshalUnMarshalPacket 结构体转字节再转结构体
func NewJsonMarshalUnMarshalPacket(dataProvider interface{}, dataReceiver interface{}) (pack stream.PacketHandlerI) {
	return &_JsonMarshalUnMarshalPacket{
		dataProvider: dataProvider,
		dataReceiver: dataReceiver,
	}
}

func (pack *_JsonMarshalUnMarshalPacket) Name() string {
	return stream.GeneratePacketHandlerName(pack)
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
	err = json.Unmarshal(input, pack.dataReceiver)
	return ctx, nil, err
}

type _JsonUnmarshalMarshalPacket struct {
	_JsonMarshalUnMarshalPacket
}

//NewJsonUnmarshalMarshalPacket 字节转结构体再转字节
func NewJsonUnmarshalMarshalPacket(dataReceiver interface{}, dataProvider interface{}) (pack stream.PacketHandlerI) {
	return &_JsonUnmarshalMarshalPacket{
		_JsonMarshalUnMarshalPacket: _JsonMarshalUnMarshalPacket{
			dataProvider: dataProvider,
			dataReceiver: dataReceiver,
		},
	}
}

func (pack *_JsonUnmarshalMarshalPacket) Name() string {
	return stream.GeneratePacketHandlerName(pack)
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
