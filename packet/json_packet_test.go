package packet_test

import (
	"fmt"
	"testing"

	"github.com/suifengpiao14/stream/packet"
)

type User struct {
	ID string
}

func TestNewJsonMarshalUnMarshalPacket(t *testing.T) {
	packetHander := packet.NewJsonMarshalUnMarshalPacket(nil, nil)
	fmt.Println(packetHander.Name())
}
