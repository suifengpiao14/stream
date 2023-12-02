package packet

import (
	"context"
	"encoding/json"

	"github.com/suifengpiao14/lineschema"
	"github.com/suifengpiao14/stream"
)

type _JsonMergePacket struct {
	BeforeMergedData []byte `json:"beforeMergedData"`
	AfterMergedData  []byte `json:"afterMergedData"`
}

// NewJsonMergePacket 合并数据
func NewJsonMergePacket(beforeMergedData []byte, afterMergedData []byte) (pack stream.PacketHandlerI) {
	return &_JsonMergePacket{
		BeforeMergedData: beforeMergedData,
		AfterMergedData:  afterMergedData,
	}
}

func (pack *_JsonMergePacket) Name() string {
	return stream.GeneratePacketHandlerName(pack)
}
func (pack *_JsonMergePacket) Description() string {
	return "合并json数据"
}

func (pack *_JsonMergePacket) String() string {
	b, _ := json.Marshal(pack)
	s := string(b)
	return s
}

func (pack *_JsonMergePacket) Before(ctx context.Context, input []byte) (newCtx context.Context, out []byte, err error) {
	if pack.BeforeMergedData == nil {
		return ctx, input, nil
	}
	out, err = lineschema.MergeDefault(input, pack.BeforeMergedData)
	if err != nil {
		return nil, nil, err
	}
	return ctx, out, nil
}

func (pack *_JsonMergePacket) After(ctx context.Context, input []byte) (newCtx context.Context, out []byte, err error) {
	if pack.AfterMergedData == nil {
		return ctx, input, nil
	}
	out, err = lineschema.MergeDefault(input, pack.AfterMergedData)
	if err != nil {
		return nil, nil, err
	}
	return ctx, out, nil
}
