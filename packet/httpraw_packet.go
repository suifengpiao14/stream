package packet

import (
	"context"
	"encoding/json"

	"github.com/suifengpiao14/httpraw"
	"github.com/suifengpiao14/packethandler"
)

type _HttprawPacketHandler struct {
}

const PACKETHANDLER_NAME_HttprawPacketHandler = "github.com/suifengpiao14/stream/packet/_HttprawPacketHandler"

func NewHttprawPacketHandler() (packHandler packethandler.PacketHandlerI) {
	return &_HttprawPacketHandler{}
}

func (packet *_HttprawPacketHandler) Name() string {
	return PACKETHANDLER_NAME_HttprawPacketHandler
}

func (packet *_HttprawPacketHandler) Description() string {
	return `将http 协议转换为预定义的json格式数据`
}
func (packet *_HttprawPacketHandler) Before(ctx context.Context, input []byte) (newCtx context.Context, out []byte, err error) {
	wellHttpRaw := string(input)
	r, err := httpraw.ReadRequest(wellHttpRaw)
	if err != nil {
		return ctx, nil, err
	}
	rDTO, err := httpraw.DestructReqeust(r)
	if err != nil {
		return ctx, nil, err
	}
	out, err = json.Marshal(rDTO)
	if err != nil {
		return ctx, nil, err
	}
	return ctx, out, nil
}
func (packet *_HttprawPacketHandler) After(ctx context.Context, input []byte) (newCtx context.Context, out []byte, err error) {
	return packethandler.EmptyHandlerFn(ctx, input)
}

func (packet *_HttprawPacketHandler) String() string {
	return ""
}
