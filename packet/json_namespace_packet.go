package packet

import (
	"context"

	"github.com/suifengpiao14/stream"
	"github.com/tidwall/gjson"
	"github.com/tidwall/sjson"
)

type _JsonAddTrimNamespacePacket struct {
	namespace string
}

// NewJsonAddTrimNamespacePacket 给json增加命名空间
func NewJsonAddTrimNamespacePacket(namespace string) (pack stream.PacketHandlerI) {
	return &_JsonAddTrimNamespacePacket{
		namespace: namespace,
	}
}

func (pack *_JsonAddTrimNamespacePacket) Name() string {
	return stream.GeneratePacketHandlerName(pack)
}
func (pack *_JsonAddTrimNamespacePacket) Description() string {
	return "add namespace to json"
}

func (pack *_JsonAddTrimNamespacePacket) String() string {
	return ""
}

func (pack *_JsonAddTrimNamespacePacket) Before(ctx context.Context, input []byte) (newCtx context.Context, out []byte, err error) {
	if pack.namespace == "" {
		return ctx, input, nil
	}
	out, err = sjson.SetRawBytes([]byte{}, pack.namespace, input)
	if err != nil {
		return nil, nil, err
	}
	return ctx, out, nil

}

func (pack *_JsonAddTrimNamespacePacket) After(ctx context.Context, input []byte) (newCtx context.Context, out []byte, err error) {
	if pack.namespace == "" {
		return ctx, input, nil
	}
	s := gjson.GetBytes(input, pack.namespace).String()
	out = []byte(s)
	return ctx, out, nil
}

type _JsonTrimAddNamespacePacket struct {
	_JsonAddTrimNamespacePacket
}

// NewJsonTrimAddNamespacePacket 删除json命名空间
func NewJsonTrimAddNamespacePacket(namespace string) (pack stream.PacketHandlerI) {
	return &_JsonTrimAddNamespacePacket{
		_JsonAddTrimNamespacePacket: _JsonAddTrimNamespacePacket{
			namespace: namespace,
		},
	}
}

func (pack *_JsonTrimAddNamespacePacket) Name() string {
	return stream.GeneratePacketHandlerName(pack)
}
func (pack *_JsonTrimAddNamespacePacket) Description() string {
	return "drop json name space"
}

func (pack *_JsonTrimAddNamespacePacket) Before(ctx context.Context, input []byte) (newCtx context.Context, out []byte, err error) {
	return pack._JsonAddTrimNamespacePacket.After(ctx, input)
}

func (pack *_JsonTrimAddNamespacePacket) After(ctx context.Context, input []byte) (newCtx context.Context, out []byte, err error) {
	return pack._JsonAddTrimNamespacePacket.Before(ctx, input)
}
