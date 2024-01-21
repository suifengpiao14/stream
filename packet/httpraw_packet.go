package packet

import (
	"context"
	"encoding/json"

	"github.com/suifengpiao14/httpraw"
	"github.com/suifengpiao14/stream"
)

type _HttprawPacketHandler struct {
}

func NewHttprawPacketHandler() (packHandler stream.PacketHandlerI) {
	return &_HttprawPacketHandler{}
}

func (packet *_HttprawPacketHandler) Name() string {
	return stream.GeneratePacketHandlerName(packet)
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
	return
}

func (packet *_HttprawPacketHandler) String() string {
	return ""
}
