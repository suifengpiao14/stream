package stream

import (
	"sync"

	"github.com/pkg/errors"
)

type Container struct {
	lock    sync.Mutex
	streams map[string]*Stream
}

func NewContainer() *Container {
	return &Container{}
}

func (c *Container) RegisterStream(stream *Stream) (err error) {
	c.lock.Lock()
	defer c.lock.Unlock()
	if _, ok := c.streams[stream.Name]; ok {
		err = errors.Errorf("stream register twice by name:%s", stream.Name)
		return err
	}
	c.streams[stream.Name] = stream
	return err
}
func (c *Container) GetStream(name string) (stream *Stream, err error) {
	stream, ok := c.streams[stream.Name]
	if !ok {
		err = errors.Errorf("stream not found by name:%s", name)
		return nil, err
	}
	return stream, nil
}
