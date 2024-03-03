package sqlpluspack

import (
	"context"
	"encoding/json"

	"github.com/suifengpiao14/sqlplus"
	"github.com/suifengpiao14/stream"
)

type SqlPlusPacketHandler struct {
	Scenes       sqlplus.Scenes        `json:"scenes"`
	TableColumns []sqlplus.TableColumn `json:"tableColumns"`
}

const PACKETHANDLER_NAME_SqlPlusPacketHandler = "github.com/suifengpiao14/stream/packet/sqlpluspack/SqlPlusPacketHandler"

// NewSqlPlusPacketHandler sql 增删改查语句扩展
func NewSqlPlusPacketHandler(scenes sqlplus.Scenes, tableColumns ...sqlplus.TableColumn) (packHandler stream.PacketHandlerI) {
	return &SqlPlusPacketHandler{
		Scenes:       scenes,
		TableColumns: tableColumns,
	}
}

func (packet *SqlPlusPacketHandler) Name() string {
	return PACKETHANDLER_NAME_SqlPlusPacketHandler
}
func (packet *SqlPlusPacketHandler) String() string {
	b, _ := json.Marshal(packet)
	s := string(b)
	return s
}

func (packet *SqlPlusPacketHandler) Description() string {
	return `根据配置,解析传入的sql,修改、新增内容后再形成新的sql输出`
}

func (packet *SqlPlusPacketHandler) Before(ctx context.Context, input []byte) (newCtx context.Context, out []byte, err error) {
	sql := string(input)
	newSql, err := sqlplus.WithPlusScene(sql, packet.Scenes, packet.TableColumns...)
	if err != nil {
		return ctx, nil, err
	}
	out = []byte(newSql)
	return ctx, out, nil
}

func (packet *SqlPlusPacketHandler) After(ctx context.Context, input []byte) (newCtx context.Context, out []byte, err error) {
	return ctx, input, nil
}
