package stream_test

import (
	"fmt"
	"testing"

	"github.com/suifengpiao14/stream"
)

type User struct {
	ID string
}

func TestNewJsonMarshalUnMarshalPacket(t *testing.T) {
	packetHander := stream.NewJsonMarshalUnMarshalPacket(nil, nil)
	fmt.Println(packetHander.Name())
}
