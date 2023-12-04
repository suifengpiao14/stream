package stream_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/suifengpiao14/stream"
	"github.com/suifengpiao14/stream/packet"
)

func TestInsertBefor(t *testing.T) {
	var dataProvider interface{}
	var dataReceiver interface{}
	marshalUnMarshal := packet.NewJsonMarshalUnMarshalPacket(&dataProvider, &dataReceiver)
	unMarshalMarshal := packet.NewJsonUnmarshalMarshalPacket(&dataReceiver, &dataProvider)
	commonPackets := stream.NewPacketHandlers(
		marshalUnMarshal,
		unMarshalMarshal,
	)
	transferPacket := packet.NewTransferPacketHandler("", "")
	t.Run("after not found", func(t *testing.T) {
		packets := commonPackets
		packets.InsertAfter(2, transferPacket)
		assert.Equal(t, transferPacket.Name(), packets[2].Name())
	})
	t.Run("after last", func(t *testing.T) {
		packets := commonPackets
		packets.InsertAfter(1, transferPacket)
		assert.Equal(t, transferPacket.Name(), packets[2].Name())
	})
	t.Run("after  first", func(t *testing.T) {
		packets := commonPackets
		packets.InsertAfter(0, transferPacket)
		assert.Equal(t, transferPacket.Name(), packets[1].Name())
	})

	t.Run("not found", func(t *testing.T) {
		packets := commonPackets
		packets.InsertBefore(-5, transferPacket)
		assert.Equal(t, transferPacket.Name(), packets[0].Name())
	})
	t.Run("befor first", func(t *testing.T) {
		packets := commonPackets
		packets.InsertBefore(0, transferPacket)
		assert.Equal(t, transferPacket.Name(), packets[0].Name())
	})
	t.Run("before  last", func(t *testing.T) {
		packets := commonPackets
		packets.InsertBefore(1, transferPacket)
		assert.Equal(t, transferPacket.Name(), packets[1].Name())
	})

	t.Run("Delete first", func(t *testing.T) {
		packets := commonPackets
		packets.Delete(0)
		assert.Equal(t, unMarshalMarshal.Name(), packets[0].Name())
	})
	t.Run("Delete last", func(t *testing.T) {
		packets := commonPackets
		packets.Delete(1)
		assert.Equal(t, marshalUnMarshal.Name(), packets[0].Name())
		assert.Equal(t, 1, len(packets))
	})

	t.Run("Delete middle", func(t *testing.T) {
		packets := commonPackets
		packets.InsertBefore(1, transferPacket)
		packets.Delete(1)
		assert.Equal(t, unMarshalMarshal.Name(), packets[1].Name())
		assert.Equal(t, 2, len(packets))
	})
	t.Run("replace first", func(t *testing.T) {
		packets := commonPackets
		newMarshalUnMarshal := packet.NewJsonMarshalUnMarshalPacket(nil, nil)
		packets.Replace(0, newMarshalUnMarshal)
		assert.Equal(t, 2, len(packets))
	})

}
