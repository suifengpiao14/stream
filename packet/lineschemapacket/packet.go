package lineschemapacket

import (
	"fmt"

	"github.com/suifengpiao14/lineschema"
	"github.com/suifengpiao14/stream"
)

// lineschema 格式数据包
type LineschemaPacketI interface {
	GetRoute() (mehtod string, path string) // 网络传输地址，http可用method,path标记
	UnpackSchema() (lineschema string)      // 解包配置 从网络数据到程序
	PackSchema() (lineschema string)        // 封包配置 程序到网络
}

func RegisterLineschemaPacket(pack LineschemaPacketI) (err error) {
	method, path := pack.GetRoute()
	unpackId, packId := makeLineschemaPacketKey(method, path)
	unpackSchema, packSchema := pack.UnpackSchema(), pack.PackSchema()
	unpackLineschema, err := lineschema.ParseLineschema(unpackSchema)
	if err != nil {
		return err
	}
	packLineschema, err := lineschema.ParseLineschema(packSchema)
	if err != nil {
		return err
	}
	err = RegisterLineschema(unpackId, *unpackLineschema)
	if err != nil {
		return err
	}
	err = RegisterLineschema(packId, *packLineschema)
	if err != nil {
		return err
	}
	return err
}

func ServerpacketHandlers(api LineschemaPacketI) (packetHandlers stream.PacketHandlers, err error) {
	method, path := api.GetRoute()
	unpackId, packId := makeLineschemaPacketKey(method, path)
	unpackLineschema, err := GetClineschema(unpackId)
	if err != nil {
		return nil, err
	}

	packLineschema, err := GetClineschema(packId)
	if err != nil {
		return nil, err
	}
	packetHandlers = make(stream.PacketHandlers, 0)
	packetHandlers.Append(
		NewValidatePacketHandler(string(unpackLineschema.Jsonschema), string(packLineschema.Jsonschema), unpackLineschema.validateLoader, packLineschema.validateLoader),
		NewMergeDefaultHandler(string(unpackLineschema.DefaultJson), string(packLineschema.DefaultJson)),
		NewTransferPacketHandler(unpackLineschema.transferToFormatGjsonPath, packLineschema.transferToTypeGjsonPath),
	)
	return packetHandlers, nil
}

func SDKpacketHandlers(api LineschemaPacketI) (packetHandlers stream.PacketHandlers, err error) {
	method, path := api.GetRoute()
	unpackId, packId := makeLineschemaPacketKey(method, path)
	unpackLineschema, err := GetClineschema(unpackId)
	if err != nil {
		return nil, err
	}

	packLineschema, err := GetClineschema(packId)
	if err != nil {
		return nil, err
	}
	packetHandlers = make(stream.PacketHandlers, 0)
	packetHandlers.Append(
		NewTransferPacketHandler(packLineschema.transferToTypeGjsonPath, unpackLineschema.transferToFormatGjsonPath),
		NewMergeDefaultHandler(string(packLineschema.DefaultJson), string(unpackLineschema.DefaultJson)),
		NewValidatePacketHandler(string(packLineschema.Jsonschema), string(unpackLineschema.Jsonschema), packLineschema.validateLoader, unpackLineschema.validateLoader),
	)
	return packetHandlers, nil
}

func makeLineschemaPacketKey(method string, path string) (unpackId string, packId string) {
	unpackId = fmt.Sprintf("%s-%s-input", method, path)
	packId = fmt.Sprintf("%s-%s-output", method, path)
	return unpackId, packId
}
