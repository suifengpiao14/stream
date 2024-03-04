package packet

import (
	"context"
	"database/sql"

	"github.com/blastrain/vitess-sqlparser/sqlparser"
	"github.com/spf13/cast"
	"github.com/suifengpiao14/cudevent/cudeventimpl"
	"github.com/suifengpiao14/packethandler"
	"github.com/suifengpiao14/sqlexec"
	"github.com/suifengpiao14/sqlplus"
)

type _CUDEventPackHandler struct {
	db          *sql.DB
	database    string
	sqlRawEvent *cudeventimpl.SQLRawEvent
}

const PACKETHANDLER_NAME_CUDEvent = "github.com/suifengpiao14/stream/packet/_CUDEventPackHandler"

func NewCUDEventPackHandler(db *sql.DB, dbName string) (packHandler packethandler.PacketHandlerI) {
	packHandler = &_CUDEventPackHandler{
		db:       db,
		database: dbName,
	}
	return packHandler
}

func (packet *_CUDEventPackHandler) Name() string {
	return PACKETHANDLER_NAME_CUDEvent
}

func (packet *_CUDEventPackHandler) Description() string {
	return `解析sql,发布增改删事件`
}

func (packet *_CUDEventPackHandler) String() string {
	return ""
}

func (packet *_CUDEventPackHandler) Before(ctx context.Context, input []byte) (newCtx context.Context, out []byte, err error) {
	sql := string(input)
	stmt, err := sqlparser.Parse(sql)
	if err != nil {
		return ctx, nil, err
	}
	packet.sqlRawEvent = &cudeventimpl.SQLRawEvent{
		SQL:      sql,
		DB:       packet.db,
		Stmt:     stmt,
		Database: packet.database,
	} // 重新初始化
	switch stmt := stmt.(type) {
	case *sqlparser.Update: // 更新类型，先查询更新前数据，并保存
		selectSQL := sqlplus.ConvertUpdateToSelect(stmt)
		before, err := sqlexec.QueryContext(ctx, packet.db, selectSQL)
		if err != nil {
			return ctx, nil, err
		}
		packet.sqlRawEvent.BeforeData = before
	}
	return ctx, input, nil
}
func (packet *_CUDEventPackHandler) After(ctx context.Context, input []byte) (newCtx context.Context, out []byte, err error) {
	stmt := packet.sqlRawEvent.Stmt
	switch stmt.(type) {
	case *sqlparser.Insert:
		packet.sqlRawEvent.LastInsertId = string(input)
	case *sqlparser.Update:
		packet.sqlRawEvent.RowsAffected = cast.ToInt64(string(input))
	}
	err = cudeventimpl.PublishSQLRawEvent(packet.sqlRawEvent)
	if err != nil {
		return ctx, nil, err
	}
	return ctx, input, nil
}
