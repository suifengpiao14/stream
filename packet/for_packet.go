package packet

import (
	"context"
	"errors"

	"github.com/suifengpiao14/packethandler"
	"github.com/suifengpiao14/stream"
)

type _ForPacketHandler struct {
	packetHandlers packethandler.PacketHandlers
}

const PACKETHANDLER_NAME_ForPacketHandler = "github.com/suifengpiao14/stream/packet/_ForPacketHandler"

var Error_break = errors.New("break for loop")

func NewForPacketHandler(packetHandlers ...packethandler.PacketHandlerI) (packHandler packethandler.PacketHandlerI) {
	return &_ForPacketHandler{
		packetHandlers: packetHandlers,
	}
}

func (packet *_ForPacketHandler) Name() string {
	return PACKETHANDLER_NAME_ForPacketHandler
}

func (packet *_ForPacketHandler) Description() string {
	return `将for循环封装成流处理器,通过返回 Error_break 类型错误推出循环`
}
func (packet *_ForPacketHandler) Before(ctx context.Context, input []byte) (newCtx context.Context, out []byte, err error) {

	for {
		select {
		case <-ctx.Done(): // 监听上下文取消
			err = ctx.Err()
			return ctx, nil, err
		default:
		}
		s := stream.NewStream("for", nil, packet.packetHandlers...)
		input, err = s.Run(ctx, input)
		if errors.Is(err, Error_break) {
			break
		}
		if err != nil {
			return ctx, nil, err
		}
	}

	return ctx, input, nil
}
func (packet *_ForPacketHandler) After(ctx context.Context, input []byte) (newCtx context.Context, out []byte, err error) {
	return ctx, input, nil
}

func (packet *_ForPacketHandler) String() string {
	return ""
}
