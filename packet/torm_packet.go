package packet

import (
	"context"
	"encoding/json"

	"github.com/suifengpiao14/packethandler"
	"github.com/suifengpiao14/torm"
)

type _TormPackHandler struct {
	torm torm.Torm
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

	var m map[string]any
	if len(input) > 0 {
		err = json.Unmarshal(input, &m)
		if err != nil {
			return ctx, nil, err
		}
	}
	ConvertFloatsToInt(m) // 修改float64
	volume := torm.VolumeMap(m)
	sqls, _, _, err := torm.GetSQLFromTemplate(packet.torm.GetRootTemplate(), packet.torm.Name, &volume)
	if err != nil {
		return ctx, nil, err
	}
	out = []byte(sqls)
	return ctx, out, nil

}

func (packet *_TormPackHandler) After(ctx context.Context, input []byte) (newCtx context.Context, out []byte, err error) {
	return ctx, input, packethandler.ERROR_EMPTY_FUNC
}

// NewTormPackHandler 执行模板返回SQL
func NewTormPackHandler(torm torm.Torm) (packHandler packethandler.PacketHandlerI) {
	return &_TormPackHandler{
		torm: torm,
	}
}

// ConvertFloatsToInt json.Unmarsha 后整数改成float64了，此处尝试优先使用int
func ConvertFloatsToInt(data map[string]interface{}) {
	for key, value := range data {
		switch v := value.(type) {
		case float64:
			// 尝试将 float64 转换为 int
			if float64(int(v)) == v {
				data[key] = int(v)
			}
		case map[string]interface{}:
			// 递归处理嵌套的 map
			ConvertFloatsToInt(v)
		}
	}
}
