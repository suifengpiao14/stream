package packet

import (
	"context"
	"strings"

	"github.com/suifengpiao14/stream"
	"github.com/tidwall/gjson"
	"github.com/tidwall/sjson"
)

type _JsonAddTrimNamespacePacket struct {
	namespaceAdd  string
	namespaceTrim string
}

const PACKETHANDLER_NAME_JsonAddTrimNamespacePacket = "github.com/suifengpiao14/stream/packet/_JsonAddTrimNamespacePacket"

// NewJsonAddTrimNamespacePacket 给json增加命名空间
func NewJsonAddTrimNamespacePacket(namespaceAdd string, namespaceTrim string) (pack stream.PacketHandlerI) {
	return &_JsonAddTrimNamespacePacket{
		namespaceAdd:  strings.TrimSuffix(namespaceAdd, "."),
		namespaceTrim: strings.TrimSuffix(namespaceTrim, "."),
	}
}

func (pack *_JsonAddTrimNamespacePacket) Name() string {
	return PACKETHANDLER_NAME_JsonAddTrimNamespacePacket
}
func (pack *_JsonAddTrimNamespacePacket) Description() string {
	return "add namespace to json"
}

func (pack *_JsonAddTrimNamespacePacket) String() string {
	return ""
}

func (pack *_JsonAddTrimNamespacePacket) Before(ctx context.Context, input []byte) (newCtx context.Context, out []byte, err error) {
	if pack.namespaceAdd == "" {
		return ctx, input, nil
	}
	out, err = sjson.SetRawBytes([]byte{}, pack.namespaceAdd, input)
	if err != nil {
		return nil, nil, err
	}
	return ctx, out, nil

}

func (pack *_JsonAddTrimNamespacePacket) After(ctx context.Context, input []byte) (newCtx context.Context, out []byte, err error) {
	if pack.namespaceTrim == "" {
		return ctx, input, nil
	}
	s := gjson.GetBytes(input, pack.namespaceTrim).String()
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
			namespaceAdd: namespace,
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
