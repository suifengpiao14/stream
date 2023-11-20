package stream

import (
	"context"
	"encoding/json"
)

// Struct2Bytes2StructJsonPacket 先将对象转换为byte，一系列处理后转换为新的对象，sdk调用场景，stream 的最外层
func Struct2Bytes2StructJsonPacket(dataProvider interface{}, dataReceiver interface{}) (pack PackHandler) {
	return NewPackHandler(
		func(ctx context.Context, input []byte) (out []byte, err error) {
			return json.Marshal(dataProvider)
		},
		func(ctx context.Context, input []byte) (out []byte, err error) {
			return nil, json.Unmarshal(input, dataReceiver)
		},
	)
}

// Bytes2Stuct2BytesJsonPacket 先将byte转换为对象，一系列处理后转换为新的byte，server服务场景，stream 的最里层
func Bytes2Stuct2BytesJsonPacket(dataReceiver interface{}, dataProvider interface{}) (pack PackHandler) {
	return NewPackHandler(
		func(ctx context.Context, input []byte) (out []byte, err error) {
			return nil, json.Unmarshal(input, dataReceiver)
		},
		func(ctx context.Context, input []byte) (out []byte, err error) {
			return json.Marshal(dataProvider)
		},
	)
}
