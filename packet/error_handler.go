package packet

import (
	"context"

	"github.com/suifengpiao14/stream"
)

type ErrorI interface {
	Error() (err error)
}

type ErrorPacketHandler struct {
	BeforeErr ErrorI
	AfterErr  ErrorI
}

const PACKETHANDLER_NAME_ErrorPacketHandler = "github.com/suifengpiao14/stream/packet/ErrorPacketHandler"

func NewErrorIPacketHandler(beforeErr ErrorI, afterErr ErrorI) (packet stream.PacketHandlerI) {
	return &ErrorPacketHandler{
		BeforeErr: beforeErr,
		AfterErr:  afterErr,
	}
}

func (packet *ErrorPacketHandler) Name() string {
	return PACKETHANDLER_NAME_ErrorPacketHandler
}
func (packet *ErrorPacketHandler) Description() string {
	return `将实现了ErrorI接口的结构体,封装成packet,目的是获取其中的error`
}

func (packet *ErrorPacketHandler) String() (s string) {
	var err error
	if packet.AfterErr != nil {
		err = packet.AfterErr.Error()
	}
	if err != nil {
		s = err.Error()
	}
	return s
}
func (packet *ErrorPacketHandler) Before(ctx context.Context, input []byte) (newCtx context.Context, out []byte, err error) {
	if packet.BeforeErr != nil {
		err = packet.AfterErr.Error()
	}
	if err != nil {
		return ctx, input, err
	}
	return ctx, input, nil
}

func (packet *ErrorPacketHandler) After(ctx context.Context, input []byte) (newCtx context.Context, out []byte, err error) {
	if packet.BeforeErr != nil {
		err = packet.AfterErr.Error()
	}
	if err != nil {
		return ctx, input, err
	}
	return ctx, input, nil
}
