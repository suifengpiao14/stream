package packet

import (
	"context"

	"github.com/suifengpiao14/stream"
)

type ErrorPacketHandler struct {
	Error error
}

func NewErrorIPacketHandler(err error) (packet stream.PacketHandlerI) {
	return &ErrorPacketHandler{
		Error: err,
	}
}

func (packet *ErrorPacketHandler) Name() string {
	return stream.GeneratePacketHandlerName(packet)
}
func (packet *ErrorPacketHandler) Description() string {
	return `将实现了ErrorI接口的结构体,封装成packet,目的是获取其中的error`
}

func (packet *ErrorPacketHandler) String() string {
	return packet.Error.Error()
}
func (packet *ErrorPacketHandler) Before(ctx context.Context, input []byte) (newCtx context.Context, out []byte, err error) {
	if packet.Error.Error() == "" {
		return ctx, input, nil
	}
	return ctx, input, packet.Error
}

func (packet *ErrorPacketHandler) After(ctx context.Context, input []byte) (newCtx context.Context, out []byte, err error) {
	if packet.Error.Error() == "" {
		return ctx, input, nil
	}
	return ctx, input, packet.Error
}
