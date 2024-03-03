package yaegipacket

import (
	"context"
)

const PACKETHANDLER_NAME_YaegiHook = "github.com/suifengpiao14/stream/packet/yaegipacket/YaegiHook"

func (packet *YaegiHook) Name() string {
	return PACKETHANDLER_NAME_YaegiHook
}
func (packet *YaegiHook) Description() string {
	return ``
}

func (packet *YaegiHook) String() string {
	return packet.dynamicScript
}
func (packet *YaegiHook) Before(ctx context.Context, input []byte) (newCtx context.Context, out []byte, err error) {
	if packet.DynamicBefore == nil {
		return ctx, input, nil
	}
	return packet.DynamicBefore(ctx, input)
}

func (packet *YaegiHook) After(ctx context.Context, input []byte) (newCtx context.Context, out []byte, err error) {
	if packet.DynamicAfter == nil {
		return ctx, input, nil
	}
	return packet.DynamicAfter(ctx, input)
}
