package packet

import (
	"context"

	"github.com/suifengpiao14/stream"
)

type _FuncPacketHandler struct {
	beforeFn stream.HandlerFn
	afterFn  stream.HandlerFn
}

func NewFuncPacketHandler(beforeFn stream.HandlerFn, afterFn stream.HandlerFn) (packHandler stream.PacketHandlerI) {
	return &_FuncPacketHandler{
		beforeFn: beforeFn,
		afterFn:  afterFn,
	}
}

func (packet *_FuncPacketHandler) Name() string {
	return stream.GeneratePacketHandlerName(packet)
}

func (packet *_FuncPacketHandler) Description() string {
	return `将函数封装为packet`
}
func (packet *_FuncPacketHandler) Before(ctx context.Context, input []byte) (newCtx context.Context, out []byte, err error) {

	if packet.afterFn == nil {
		return ctx, input, nil
	}
	newCtx, out, err = packet.beforeFn(ctx, input)
	if err != nil {
		return ctx, nil, err
	}
	return newCtx, out, nil
}
func (packet *_FuncPacketHandler) After(ctx context.Context, input []byte) (newCtx context.Context, out []byte, err error) {
	if packet.afterFn == nil {
		return ctx, input, nil
	}
	newCtx, out, err = packet.afterFn(ctx, input)
	if err != nil {
		return ctx, nil, err
	}
	return newCtx, out, nil
}

func (packet *_FuncPacketHandler) String() string {
	return ""
}
