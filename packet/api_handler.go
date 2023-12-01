package packet

import (
	"context"

	"github.com/suifengpiao14/apihandler"
	"github.com/suifengpiao14/stream"
)

type _ApiPackHandler struct {
	api apihandler.ApiInterface
}

func NewApiPackHandler(api apihandler.ApiInterface) (packet stream.PacketHandlerI) {
	return &_ApiPackHandler{
		api: api,
	}
}

func (packet *_ApiPackHandler) Name() string {
	return stream.GeneratePacketHandlerName(packet)
}

func (packet *_ApiPackHandler) Description() string {
	return `封装Api接口`
}

func (packet *_ApiPackHandler) String() string {
	return ""
}

func (packet *_ApiPackHandler) Before(ctx context.Context, input []byte) (newCtx context.Context, out []byte, err error) {
	err = packet.api.Do(ctx)
	if err != nil {
		return ctx, nil, err
	}
	return ctx, nil, nil
}

func (packet *_ApiPackHandler) After(ctx context.Context, input []byte) (newCtx context.Context, out []byte, err error) {
	return ctx, input, nil
}

func ApiPackHandlers(api apihandler.ApiInterface) (packHandlers stream.PackHandlers) {
	packHandlers = make(stream.PackHandlers, 0)
	packHandlers.Add(
		NewJsonUnmarshalMarshalPacket(api, api.GetOutRef()),
		NewApiPackHandler(api),
	)
	return packHandlers
}
