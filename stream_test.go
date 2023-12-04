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
	packets := stream.NewPacketHandlers(
		marshalUnMarshal,
		unMarshalMarshal,
	)
	transferPacket := packet.NewTransferPacketHandler("", "")
	t.Run("after not found", func(t *testing.T) {
		packets.InsertAfter("", transferPacket)
		assert.Equal(t, transferPacket.Name(), packets[2].Name())
	})
	t.Run("after last", func(t *testing.T) {
		packets.InsertAfter(unMarshalMarshal.Name(), transferPacket)
		assert.Equal(t, transferPacket.Name(), packets[2].Name())
	})
	t.Run("after  first", func(t *testing.T) {
		packets.InsertAfter(marshalUnMarshal.Name(), transferPacket)
		assert.Equal(t, transferPacket.Name(), packets[1].Name())
	})

	t.Run("not found", func(t *testing.T) {
		packets.InsertBefore("", transferPacket)
		assert.Equal(t, transferPacket.Name(), packets[0].Name())
	})
	t.Run("befor first", func(t *testing.T) {
		packets.InsertBefore(marshalUnMarshal.Name(), transferPacket)
		assert.Equal(t, transferPacket.Name(), packets[0].Name())
	})
	t.Run("before  last", func(t *testing.T) {
		packets.InsertBefore(unMarshalMarshal.Name(), transferPacket)
		assert.Equal(t, transferPacket.Name(), packets[1].Name())
	})

}
