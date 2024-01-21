package packet

import (
	"bytes"
	"context"
	"encoding/json"
	"reflect"
	"text/template"

	"github.com/suifengpiao14/stream"
)

type _TemplatePacketHandler struct {
	tpl      *template.Template
	dataType reflect.Type
}

func NewTemplatePacketHandler(tpl template.Template, dataType reflect.Type) (packHandler stream.PacketHandlerI) {
	return &_TemplatePacketHandler{
		tpl:      &tpl,
		dataType: dataType,
	}
}

func (packet *_TemplatePacketHandler) Name() string {
	return stream.GeneratePacketHandlerName(packet)
}

func (packet *_TemplatePacketHandler) Description() string {
	return `执行模板转换`
}
func (packet *_TemplatePacketHandler) Before(ctx context.Context, input []byte) (newCtx context.Context, out []byte, err error) {

	var b bytes.Buffer
	zeroValue := reflect.Zero(packet.dataType)
	data := zeroValue.Interface()
	err = json.Unmarshal(input, &data)
	if err != nil {
		return ctx, nil, err
	}

	err = packet.tpl.Execute(&b, data)
	if err != nil {
		return
	}
	out = b.Bytes()
	return ctx, out, nil
}
func (packet *_TemplatePacketHandler) After(ctx context.Context, input []byte) (newCtx context.Context, out []byte, err error) {
	return
}

func (packet *_TemplatePacketHandler) String() string {
	return ""
}
