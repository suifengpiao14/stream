package packet

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/pkg/errors"
	"github.com/suifengpiao14/packethandler"
	"github.com/suifengpiao14/pathtransfer"
	"github.com/suifengpiao14/sqlexec"
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
	sqls, _, _, err := torm.GetSQLFromTemplate(packet.torm.GetRootTemplate(), packet.torm.TplName, &volume)
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

//TormSQLDefaultPacketHandler sql torm 默认处理器
func TormSQLPacketHandler(torm torm.Torm) (packetHandlers packethandler.PacketHandlers, err error) {
	packetHandlers = make(packethandler.PacketHandlers, 0)
	tormName := torm.Name()
	inputPathTransfers, outputPathTransfers := torm.Transfers.GetByNamespace(tormName).SplitInOut()
	namespaceInput := fmt.Sprintf("%s%s", tormName, pathtransfer.Transfer_Direction_input)   //去除命名空间
	namespaceOutput := fmt.Sprintf("%s%s", tormName, pathtransfer.Transfer_Direction_output) // 补充命名空间
	inputGopath := inputPathTransfers.Reverse().ModifyDstPath(func(path string) (newPath string) {
		newPath = strings.TrimPrefix(path, namespaceInput)
		return newPath
	}).GjsonPath()
	outputGopath := outputPathTransfers.ModifySrcPath(func(path string) (newPath string) {
		newPath = strings.TrimPrefix(path, namespaceOutput)
		return newPath
	}).GjsonPath()
	//转换为代码中期望的数据格式
	transferHandler := NewTransferPacketHandler(inputGopath, outputGopath)
	packetHandlers.Append(transferHandler)
	packetHandlers.Append(NewTormPackHandler(torm))

	prov := torm.Source.Provider
	dbProvider, ok := prov.(*sqlexec.ExecutorSQL)
	if !ok {
		err = errors.Errorf("ExecSQLTPL required sourceprovider.DBProvider source,got:%s", prov.TypeName())
		return nil, err
	}
	db := dbProvider.GetDB()
	if db == nil {
		err = errors.Errorf("ExecSQLTPL sourceprovider.DBProvider.GetDB required,got nil (%s)", prov.TypeName())
		return nil, err
	}
	databaseName, err := sqlexec.GetDatabaseName(db)
	if err != nil {
		return nil, err
	}
	cudeventPack := NewCUDEventPackHandler(db, databaseName)
	packetHandlers.Append(cudeventPack)
	mysqlPack := NewMysqlPacketHandler(db)
	packetHandlers.Append(mysqlPack)
	return packetHandlers, nil
}
