package packet

import (
	"github.com/suifengpiao14/sdkgolib"
	"github.com/suifengpiao14/stream"
)

func SDKPacketHandlers(client sdkgolib.ClientInterface) (packetHandlers stream.PacketHandlers) {
	packetHandlers = make(stream.PacketHandlers, 0)
	out := client.GetOutRef()
	strucpackHandler := NewJsonMarshalUnMarshalPacket(client, out)
	packetHandlers.Add(strucpackHandler)
	packetHandlers.Add(NewErrorIPacketHandler(out))
	return packetHandlers
}
