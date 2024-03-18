package packet

import (
	"context"
	"encoding/json"

	jsonpatch "github.com/evanphx/json-patch/v5"
	"github.com/suifengpiao14/packethandler"
)

type _JsonMergeInputPacket struct {
	input []byte
}

const PACKETHANDLER_NAME_JsonMergeInputPacket = "github.com/suifengpiao14/stream/packet/_JsonMergeInputPacket"

// NewJsonMergeInputPacket 输出数据中合并输入数据
func NewJsonMergeInputPacket() (pack packethandler.PacketHandlerI) {
	return &_JsonMergeInputPacket{}
}

func (pack *_JsonMergeInputPacket) Name() string {
	return PACKETHANDLER_NAME_JsonMergeInputPacket
}
func (pack *_JsonMergeInputPacket) Description() string {
	return "输出数据前,合并输入数据"
}

func (pack *_JsonMergeInputPacket) String() string {
	b, _ := json.Marshal(pack)
	s := string(b)
	return s
}

func (pack *_JsonMergeInputPacket) Before(ctx context.Context, input []byte) (newCtx context.Context, out []byte, err error) {
	pack.input = input
	return ctx, input, nil
}

func (pack *_JsonMergeInputPacket) After(ctx context.Context, input []byte) (newCtx context.Context, out []byte, err error) {
	out, err = jsonpatch.MergePatch(pack.input, input)
	if err != nil {
		return ctx, nil, err
	}
	return ctx, out, nil
}
