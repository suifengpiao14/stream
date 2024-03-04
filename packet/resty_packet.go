package packet

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/pkg/errors"
	"github.com/suifengpiao14/httpraw"
	"github.com/suifengpiao14/packethandler"
)

type _RestyPacketHandler struct {
	transport *http.Transport
}

const PACKETHANDLER_NAME_RestyPacketHandler = "github.com/suifengpiao14/stream/packet/_RestyPacketHandler"

func NewRestyPacketHandler(transport *http.Transport) (packHandler packethandler.PacketHandlerI) {
	return &_RestyPacketHandler{
		transport: transport,
	}
}

func (packet *_RestyPacketHandler) Name() string {
	return PACKETHANDLER_NAME_RestyPacketHandler
}

func (packet *_RestyPacketHandler) Description() string {
	return `resty http 请求`
}
func (packet *_RestyPacketHandler) Before(ctx context.Context, input []byte) (newCtx context.Context, out []byte, err error) {
	var reqDTo httpraw.RequestDTO
	err = json.Unmarshal(input, &reqDTo)
	if err != nil {
		err = errors.WithMessagef(err, "got input :%s ", string(input))
		return ctx, nil, err
	}
	r, err := httpraw.BuildRequest(&reqDTo)
	if err != nil {
		return nil, nil, err
	}

	out, err = httpraw.RestyRequestFn(ctx, r, packet.transport)
	if err != nil {
		return nil, nil, err
	}
	return ctx, out, nil
}
func (packet *_RestyPacketHandler) After(ctx context.Context, input []byte) (newCtx context.Context, out []byte, err error) {
	return packethandler.EmptyHandlerFn(ctx, input)
}

func (packet *_RestyPacketHandler) String() string {
	return ""
}
