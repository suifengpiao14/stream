package packet

import (
	"context"

	"github.com/suifengpiao14/sdkgolib"
	"github.com/suifengpiao14/stream"
)

type ErrorI interface {
	Error() (err error)
}

type ErrorIPacketHandler struct {
	Error ErrorI
}

func NewErrorIPacketHandler(errorI ErrorI) (packet stream.PacketHandlerI) {
	return &ErrorIPacketHandler{
		Error: errorI,
	}
}

func (packet *ErrorIPacketHandler) Name() string {
	return stream.GeneratePacketHandlerName(packet)
}
func (packet *ErrorIPacketHandler) Description() string {
	return `将实现了ErrorI接口的结构体,封装成packet,目的是获取其中的error`
}

func (packet *ErrorIPacketHandler) String() string {
	return packet.Error.Error().Error()
}
func (packet *ErrorIPacketHandler) Before(ctx context.Context, input []byte) (newCtx context.Context, out []byte, err error) {
	if packet.Error != nil && packet.Error.Error() != nil {
		err = packet.Error.Error()
	}
	return ctx, input, err
}

func (packet *ErrorIPacketHandler) After(ctx context.Context, input []byte) (newCtx context.Context, out []byte, err error) {
	if packet.Error != nil && packet.Error.Error() != nil {
		err = packet.Error.Error()
	}
	return ctx, input, err
}

func SDKPacketHandlers(client sdkgolib.ClientInterface) (packetHandlers stream.PacketHandlers) {
	packetHandlers = make(stream.PacketHandlers, 0)
	out := client.GetOutRef()
	strucpackHandler := NewJsonMarshalUnMarshalPacket(client, out)
	packetHandlers.Add(strucpackHandler)
	packetHandlers.Add(NewErrorIPacketHandler(out))
	return packetHandlers
}
