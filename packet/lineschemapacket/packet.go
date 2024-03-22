package lineschemapacket

import (
	"fmt"

	"github.com/suifengpiao14/lineschema"
	"github.com/suifengpiao14/packethandler"
)

// lineschema 格式数据包
type LineschemaPacketI interface {
	GetRoute() (mehtod string, path string) // 网络传输地址，http可用method,path标记
	UnpackSchema() (lineschema string)      // 解包配置 从网络数据到程序
	PackSchema() (lineschema string)        // 封包配置 程序到网络
}

func RegisterClineschemaPacket(pack LineschemaPacketI) (err error) {
	unpackCLineschema, packCLineschema, err := ParserLineschemaPacket2Clineschema(pack)
	if err != nil {
		return err
	}
	err = RegisterClineschema(*unpackCLineschema)
	if err != nil {
		return err
	}
	err = RegisterClineschema(*packCLineschema)
	if err != nil {
		return err
	}
	return err
}

func ParserLineschemaPacket2Clineschema(pack LineschemaPacketI) (unpackCLineschema *Clineschema, packCLineschema *Clineschema, err error) {
	method, path := pack.GetRoute()
	unpackId, packId := makeLineschemaPacketKey(method, path)
	unpackSchema, packSchema := pack.UnpackSchema(), pack.PackSchema()
	unpackLineschema, err := lineschema.ParseLineschema(unpackSchema)
	if err != nil {
		return nil, nil, err
	}
	unpackCLineschema, err = NewClineschame(unpackId, *unpackLineschema)
	if err != nil {
		return nil, nil, err
	}
	packLineschema, err := lineschema.ParseLineschema(packSchema)
	if err != nil {
		return nil, nil, err
	}
	packCLineschema, err = NewClineschame(packId, *packLineschema)
	if err != nil {
		return nil, nil, err
	}
	return unpackCLineschema, packCLineschema, nil

}

func ServerpacketHandlers(requestClineschema Clineschema, responseClineschema Clineschema) (packetHandlers packethandler.PacketHandlers) {
	packetHandlers = make(packethandler.PacketHandlers, 0)
	packetHandlers.Append(
		NewValidatePacketHandler(string(requestClineschema.Jsonschema), string(responseClineschema.Jsonschema), requestClineschema.validateLoader, responseClineschema.validateLoader),
		NewMergeDefaultHandler(string(requestClineschema.DefaultJson), string(responseClineschema.DefaultJson)),
		NewTransferPacketHandler(requestClineschema.transferToFormatGjsonPath, responseClineschema.transferToTypeGjsonPath),
	)
	return packetHandlers
}

func SDKpacketHandlers(requestClineschema Clineschema, responseClineschema Clineschema) (packetHandlers packethandler.PacketHandlers) {
	packetHandlers = make(packethandler.PacketHandlers, 0)
	packetHandlers.Append(
		NewTransferPacketHandler(responseClineschema.transferToTypeGjsonPath, responseClineschema.transferToFormatGjsonPath),
		NewMergeDefaultHandler(string(responseClineschema.DefaultJson), string(responseClineschema.DefaultJson)),
		NewValidatePacketHandler(string(responseClineschema.Jsonschema), string(responseClineschema.Jsonschema), responseClineschema.validateLoader, responseClineschema.validateLoader),
	)
	return packetHandlers
}

func makeLineschemaPacketKey(method string, path string) (unpackId string, packId string) {
	unpackId = fmt.Sprintf("%s-%s-input", method, path)
	packId = fmt.Sprintf("%s-%s-output", method, path)
	return unpackId, packId
}
