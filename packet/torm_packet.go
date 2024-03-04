package packet

import (
	"context"
	"encoding/json"

	"github.com/suifengpiao14/packethandler"
	"github.com/suifengpiao14/torm"
)

type _TormPackHandler struct {
	torm torm.TormI
}

const PACKETHANDLER_NAME_TormPackHandler = "github.com/suifengpiao14/stream/packet/_TormPackHandler"

func (packet *_TormPackHandler) Name() string {
	return PACKETHANDLER_NAME_TormPackHandler
}
func (packet *_TormPackHandler) Description() string {
	return `使用go template 生成sql语句`
}

func (packet *_TormPackHandler) String() string {
	return ""
}

func (packet *_TormPackHandler) Before(ctx context.Context, input []byte) (newCtx context.Context, out []byte, err error) {

	volume := make(torm.VolumeMap)
	err = json.Unmarshal(input, &volume)
	if err != nil {
		return ctx, nil, err
	}
	sqls, _, _, err := torm.GetSQL(packet.torm.Identity(), packet.torm.TplName(), &volume)
	if err != nil {
		return ctx, nil, err
	}
	out = []byte(sqls)
	return ctx, out, nil

}

func (packet *_TormPackHandler) After(ctx context.Context, input []byte) (newCtx context.Context, out []byte, err error) {
	return ctx, input, err
}

// NewTormPackHandler 执行模板返回SQL
func NewTormPackHandler(torm torm.TormI) (packHandler packethandler.PacketHandlerI) {
	return &_TormPackHandler{
		torm: torm,
	}
}
