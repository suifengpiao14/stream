package stream

import (
	"context"
	"encoding/json"
)

type SerializeI interface {
	Serialize(ctx context.Context, src interface{}) (out []byte, err error) // 对象序列化成字节
	Unserialize(ctx context.Context, dst interface{}) (err error)           // 字节反序列化成对象
}

//BytesToStructJsonPacket 先将对象转换为byte，一系列处理后转换为新的对象，sdk调用场景，stream 的最外层
func BytesToStructJsonPacket(src interface{}, dst interface{}) (pack PackHandler) {
	return PackHandler{
		Befor: func(ctx context.Context, _ []byte) (out []byte, err error) {
			return json.Marshal(src)
		},
		After: func(ctx context.Context, input []byte) (_ []byte, err error) {
			return nil, json.Unmarshal(input, dst)
		},
	}
}

//StructToBytesJsonPacket 先将byte转换为对象，一系列处理后转换为新的byte，server服务场景，stream 的最里层
func StructToBytesJsonPacket(src interface{}, dst interface{}) (pack PackHandler) {
	return PackHandler{
		Befor: func(ctx context.Context, input []byte) (_ []byte, err error) {
			return nil, json.Unmarshal(input, dst)
		},
		After: func(ctx context.Context, _ []byte) (out []byte, err error) {
			return json.Marshal(dst)
		},
	}
}
