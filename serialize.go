package stream

import (
	"context"
	"encoding/json"
)

//ToStructJsonPacket 先将对象转换为byte，一系列处理后转换为新的对象，sdk调用场景，stream 的最外层
func ToStructJsonPacket(src interface{}, dst interface{}) (pack PackHandler) {
	return NewPackHandler(
		func(ctx context.Context, input []byte) (out []byte, err error) {
			return json.Marshal(src)
		},
		func(ctx context.Context, input []byte) (out []byte, err error) {
			return nil, json.Unmarshal(input, dst)
		},
	)
}

//ToBytesJsonPacket 先将byte转换为对象，一系列处理后转换为新的byte，server服务场景，stream 的最里层
func ToBytesJsonPacket(src interface{}, dst interface{}) (pack PackHandler) {
	return ToStructJsonPacket(src, dst).Reverse()
}
