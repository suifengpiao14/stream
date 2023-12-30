package packet

import (
	"context"
	"encoding/json"

	"github.com/pkg/errors"
	"github.com/suifengpiao14/stream"
	"github.com/tidwall/gjson"
)

type _TransferPacketHandler struct {
	BeforGjsonPath string `json:"beforGjsonPath"`
	AfterGjsonPath string `json:"afterGjsonPath"`
}

// NewTransferPacketHandler json转换
func NewTransferPacketHandler(beforGjsonPath string, afterGjsonPath string) (packet stream.PacketHandlerI) {
	return &_TransferPacketHandler{
		BeforGjsonPath: beforGjsonPath,
		AfterGjsonPath: afterGjsonPath,
	}
}

func (packet *_TransferPacketHandler) Name() string {
	return stream.GeneratePacketHandlerName(packet)
}

func (packet *_TransferPacketHandler) Description() string {
	return `将json数据进行变换,支持类型和路径变换`
}

func (packet *_TransferPacketHandler) String() string {
	b, _ := json.Marshal(packet)
	s := string(b)
	return s
}

func (packet *_TransferPacketHandler) Before(ctx context.Context, input []byte) (newCtx context.Context, out []byte, err error) {
	if len(input) == 0 { //空数据(如不需要入参)则直接返回
		return ctx, input, nil
	}
	if packet.BeforGjsonPath == "" {
		return ctx, input, nil
	}
	if !gjson.ValidBytes(input) {
		err = errors.Errorf("input require json string ,got:%s", string(input))
		return ctx, nil, err
	}
	str := gjson.GetBytes(input, packet.BeforGjsonPath).String()
	out = []byte(str)
	return ctx, out, nil
}

func (packet *_TransferPacketHandler) After(ctx context.Context, input []byte) (newCtx context.Context, out []byte, err error) {
	if len(input) == 0 { //空数据(如数据库查不到数据)则直接返回
		return ctx, input, nil
	}
	if packet.AfterGjsonPath == "" {
		return ctx, input, nil
	}
	if !gjson.ValidBytes(input) {
		err = errors.Errorf("input require json string ,got:%s", string(input))
		return ctx, nil, err
	}
	str := gjson.GetBytes(input, packet.AfterGjsonPath).String()
	out = []byte(str)
	return ctx, out, nil
}
