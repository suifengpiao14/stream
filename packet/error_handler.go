package packet

import (
	"context"

	"github.com/suifengpiao14/stream"
)

type ErrorPacketHandler struct {
	BeforeErr error
	AfterErr  error
}

func NewErrorIPacketHandler(beforeErr error, afterErr error) (packet stream.PacketHandlerI) {
	return &ErrorPacketHandler{
		BeforeErr: beforeErr,
		AfterErr:  afterErr,
	}
}

func (packet *ErrorPacketHandler) Name() string {
	return stream.GeneratePacketHandlerName(packet)
}
func (packet *ErrorPacketHandler) Description() string {
	return `将实现了ErrorI接口的结构体,封装成packet,目的是获取其中的error`
}

func (packet *ErrorPacketHandler) String() string {
	return packet.AfterErr.Error()
}
func (packet *ErrorPacketHandler) Before(ctx context.Context, input []byte) (newCtx context.Context, out []byte, err error) {
	if packet.BeforeErr != nil && packet.BeforeErr.Error() == "" {
		return ctx, input, nil
	}
	return ctx, input, packet.BeforeErr
}

func (packet *ErrorPacketHandler) After(ctx context.Context, input []byte) (newCtx context.Context, out []byte, err error) {
	if packet.AfterErr != nil && packet.AfterErr.Error() == "" {
		return ctx, input, nil
	}
	return ctx, input, packet.AfterErr
}
