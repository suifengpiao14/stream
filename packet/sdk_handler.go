package packet

import (
	"context"

	"github.com/suifengpiao14/lineschema"
	"github.com/suifengpiao14/sdkgolib"
	"github.com/suifengpiao14/stream"
)

type _SDKPackHandler struct {
	sdk sdkgolib.ClientInterface
}

func NewSDKPackHandler(sdk sdkgolib.ClientInterface) (packet stream.PacketHandlerI) {
	return &_SDKPackHandler{
		sdk: sdk,
	}
}

func (packet *_SDKPackHandler) Name() string {
	return stream.GeneratePacketHandlerName(packet)
}

func (packet *_SDKPackHandler) Description() string {
	return `封装Api接口`
}

func (packet *_SDKPackHandler) String() string {
	return ""
}

func (packet *_SDKPackHandler) Before(ctx context.Context, input []byte) (newCtx context.Context, out []byte, err error) {
	out, err = packet.sdk.RequestHandler(ctx, input)
	if err != nil {
		return ctx, nil, err
	}
	return ctx, out, nil
}

func (packet *_SDKPackHandler) After(ctx context.Context, input []byte) (newCtx context.Context, out []byte, err error) {
	return ctx, input, nil
}

func SDKPacketJsonHandlers(client sdkgolib.ClientInterface) (packetHandlers stream.PacketHandlers) {
	packetHandlers = make(stream.PacketHandlers, 0)
	out := client.GetOutRef()
	packetHandlers.Append(NewErrorIPacketHandler(nil, out))

	strucpackHandler := NewJsonMarshalUnMarshalPacket(client, out)
	packetHandlers.Append(strucpackHandler)

	convertGpath := lineschema.ToGoTypeTransfer(out).String()
	transferPack := NewTransferPacketHandler("", convertGpath)
	packetHandlers.Append(transferPack)
	cfigStr := client.GetSDKConfig().String()
	packetHandlers.Append(NewJsonMergePacket([]byte(cfigStr), nil))
	packetHandlers.Append(NewSDKPackHandler(client))

	return packetHandlers
}
