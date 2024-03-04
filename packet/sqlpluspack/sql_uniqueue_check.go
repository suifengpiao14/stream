package sqlpluspack

import (
	"context"
	"database/sql"
	"encoding/json"

	"github.com/blastrain/vitess-sqlparser/sqlparser"
	"github.com/pkg/errors"
	"github.com/suifengpiao14/packethandler"
	"github.com/suifengpiao14/sqlexec"
	"github.com/suifengpiao14/sqlplus"
)

type SqlUniqueueCheckHandler struct {
	Database string `json:"string"`
	db       *sql.DB
}

const PACKETHANDLER_NAME_SqlUniqueueCheckHandler = "github.com/suifengpiao14/stream/packet/sqlpluspack/SqlUniqueueCheckHandler"

// NewSqlUniqueueCheckHandler sql 新增场景检测唯一键是否存在
func NewSqlUniqueueCheckHandler(database string, db *sql.DB) (packHandler packethandler.PacketHandlerI) {
	return &SqlUniqueueCheckHandler{
		Database: database,
		db:       db,
	}
}

func (packet *SqlUniqueueCheckHandler) Name() string {
	return PACKETHANDLER_NAME_SqlUniqueueCheckHandler
}
func (packet *SqlUniqueueCheckHandler) String() string {
	b, _ := json.Marshal(packet)
	s := string(b)
	return s
}

func (packet *SqlUniqueueCheckHandler) Description() string {
	return `根据ddl配置,解析传入的sql,生成唯一键查询语句,检测当前值是否已经存在`
}

var Error_Uniqueue_exists = errors.New("uniqueue exists")

func (packet *SqlUniqueueCheckHandler) Before(ctx context.Context, input []byte) (newCtx context.Context, out []byte, err error) {
	sql := string(input)
	selectSqls, err := sqlplus.WithCheckUniqueue(packet.Database, sql)
	if err != nil {
		return ctx, nil, err
	}
	for _, selectSql := range selectSqls {
		data, err := sqlexec.ExecOrQueryContext(ctx, packet.db, selectSql.Sql)
		if err != nil {
			return ctx, nil, err
		}
		if data != "" {
			err = errors.WithMessage(Error_Uniqueue_exists, sqlparser.String(selectSql.Where.WhereAndExpr()))

			return ctx, nil, err
		}
	}

	out = input
	return ctx, out, nil
}

func (packet *SqlUniqueueCheckHandler) After(ctx context.Context, input []byte) (newCtx context.Context, out []byte, err error) {
	return ctx, input, nil
}
