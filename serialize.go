package stream

import (
	"context"
	"encoding/json"
)

// JsonMarshalUnMarshalPacket befor marshal,after Unmarshal，sdk调用场景，stream 的最外层,支持仅执行before 或者after
func JsonMarshalUnMarshalPacket(dataProvider interface{}, dataReceiver interface{}) (pack PackHandler) {
	return NewPackHandler(
		func(ctx context.Context, input []byte) (out []byte, err error) {
			if dataProvider == nil {
				return input, nil
			}
			return json.Marshal(dataProvider)
		},
		func(ctx context.Context, input []byte) (out []byte, err error) {
			if dataReceiver == nil {
				return input, nil
			}
			return nil, json.Unmarshal(input, dataReceiver)
		},
	)
}

// JsonUnmarshalMarshalPacket before Unmarshal,after Marshal，server服务场景，stream 的最里层
func JsonUnmarshalMarshalPacket(dataReceiver interface{}, dataProvider interface{}) (pack PackHandler) {
	return NewPackHandler(
		func(ctx context.Context, input []byte) (out []byte, err error) {
			if dataReceiver == nil {
				return input, nil
			}
			return nil, json.Unmarshal(input, dataReceiver)
		},
		func(ctx context.Context, input []byte) (out []byte, err error) {
			if dataProvider == nil {
				return input, nil
			}
			return json.Marshal(dataProvider)
		},
	)
}
