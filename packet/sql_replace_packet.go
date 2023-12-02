package packet

import (
	"context"
	"database/sql"

	"github.com/blastrain/vitess-sqlparser/sqlparser"
	"github.com/spf13/cast"
	"github.com/suifengpiao14/cudevent/cudeventimpl"
	"github.com/suifengpiao14/sqlexec"
	"github.com/suifengpiao14/sqlplus"
	"github.com/suifengpiao14/stream"
)

type _SQLReplacePacketHandler struct {
	db          *sql.DB
	sqlRawEvent *cudeventimpl.SQLRawEvent
}

func NewSQLSetPacketHandler(db *sql.DB) (packHandler stream.PacketHandlerI) {
	return &_SQLReplacePacketHandler{
		db: db,
	}
}

func (packet *_SQLReplacePacketHandler) Name() string {
	return stream.GeneratePacketHandlerName(packet)
}

func (packet *_SQLReplacePacketHandler) Description() string {
	return `从更新语句中获取查询条件,存在则执行更新,不存在则转为insert语句,实现set功能`
}

func (packet *_SQLReplacePacketHandler) String() string {
	return ""
}

func (packet *_SQLReplacePacketHandler) Before(ctx context.Context, input []byte) (newCtx context.Context, out []byte, err error) {
	sql := string(input)
	stmt, err := sqlparser.Parse(sql)
	if err != nil {
		return ctx, nil, err
	}
	packet.sqlRawEvent = &cudeventimpl.SQLRawEvent{} // 重新初始化
	packet.sqlRawEvent.SQL = sql
	packet.sqlRawEvent.DB = packet.db
	packet.sqlRawEvent.Stmt = stmt
	switch stmt := stmt.(type) {
	case *sqlparser.Update: // 更新类型，先查询更新前数据，并保存
		selectSQL := sqlplus.ConvertUpdateToSelect(stmt)
		before, err := sqlexec.QueryContext(ctx, packet.db, selectSQL)
		if err != nil {
			return ctx, nil, err
		}
		packet.sqlRawEvent.BeforeData = before
		if before == "" { //不存在,则生成insert语句
			insertSql := sqlplus.ConvertUpdateToInsert(stmt)
			//替换为insert语句后,重新设置事件内容
			input = []byte(insertSql)
			sql := string(input)
			stmt, err := sqlparser.Parse(sql)
			if err != nil {
				return ctx, nil, err
			}
			packet.sqlRawEvent = &cudeventimpl.SQLRawEvent{} // 重新初始化
			packet.sqlRawEvent.SQL = sql
			packet.sqlRawEvent.DB = packet.db
			packet.sqlRawEvent.Stmt = stmt

		}
	}
	return ctx, input, nil
}
func (packet *_SQLReplacePacketHandler) After(ctx context.Context, input []byte) (newCtx context.Context, out []byte, err error) {
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
