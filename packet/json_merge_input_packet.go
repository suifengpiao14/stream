package packet

import (
	"context"

	jsonpatch "github.com/evanphx/json-patch/v5"
	"github.com/suifengpiao14/packethandler"
)

type _JsonMergeInputToOutputPacket struct {
	input []byte
}

const PACKETHANDLER_NAME_JsonMergeInputToOutputPacket = "github.com/suifengpiao14/stream/packet/_JsonMergeInputToOutputPacket"

// NewJsonMergeInputPacket 输出数据中合并输入数据
func NewJsonMergeInputPacket() (pack packethandler.PacketHandlerI) {
	return &_JsonMergeInputToOutputPacket{}
}

func (pack *_JsonMergeInputToOutputPacket) Name() string {
	return PACKETHANDLER_NAME_JsonMergeInputToOutputPacket
}
func (pack *_JsonMergeInputToOutputPacket) Description() string {
	return "输出数据前,合并输入数据"
}

func (pack *_JsonMergeInputToOutputPacket) String() string {
	s := string(pack.input)
	return s
}

func (pack *_JsonMergeInputToOutputPacket) Before(ctx context.Context, input []byte) (newCtx context.Context, out []byte, err error) {
	pack.input = input
	return ctx, input, nil
}

func (pack *_JsonMergeInputToOutputPacket) After(ctx context.Context, input []byte) (newCtx context.Context, out []byte, err error) {
	if pack.input == nil {
		return ctx, input, nil
	}

	if input == nil {
		return ctx, pack.input, nil
	}
	out, err = jsonpatch.MergePatch(pack.input, input)
	if err != nil {
		return ctx, nil, err
	}
	return ctx, out, nil
}
