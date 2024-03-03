package stream

import (
	"context"
	"encoding/json"
	"strings"

	"github.com/pkg/errors"
)

var ERROR_EMPTY_FUNC = errors.New("empty func")

// 定回调函数指针的类型
type HandlerFn func(ctx context.Context, input []byte) (newCtx context.Context, out []byte, err error)
type ErrorHandler func(ctx context.Context, err error) (out []byte)

func EmptyHandlerFn(ctx context.Context, input []byte) (newCtx context.Context, out []byte, err error) {
	return ctx, input, ERROR_EMPTY_FUNC
}

type PacketHandlerI interface {
	Name() string
	Description() string
	Before(ctx context.Context, input []byte) (newCtx context.Context, out []byte, err error)
	After(ctx context.Context, input []byte) (newCtx context.Context, out []byte, err error)
	String() string
}

// JsonString 用户实现PacketHandlerI.String()
func JsonString(packet PacketHandlerI) string {
	b, _ := json.Marshal(packet)
	s := string(b)
	return s
}

type PacketHandlers []PacketHandlerI

func NewPacketHandlers(packetHandlerIs ...PacketHandlerI) (packetHandlers PacketHandlers) {
	packetHandlers = make(PacketHandlers, 0)
	packetHandlers.Append(packetHandlerIs...)
	return packetHandlers
}

func (ps *PacketHandlers) Append(packetHandlers ...PacketHandlerI) {
	if *ps == nil {
		*ps = make(PacketHandlers, 0)
	}
	*ps = append(*ps, packetHandlers...)
}

// GetByName 通过名称获取子集合
func (ps *PacketHandlers) GetNames() (names []string) {
	names = make([]string, 0)
	for _, p := range *ps {
		names = append(names, p.Name())
	}
	return names
}

// GetByName 通过名称获取子集合
func (ps *PacketHandlers) GetByName(names ...string) (packetHandlers PacketHandlers, err error) {
	packetHandlers = make(PacketHandlers, 0)
	for _, name := range names {
		exists := false
		for _, p := range *ps {
			if strings.EqualFold(p.Name(), name) {
				packetHandlers.Append(p)
				exists = true
				break
			}
		}
		if !exists {
			err = errors.Errorf("not found packet handler named:%s", name)
			return nil, err
		}
	}

	return packetHandlers, nil
}

func (ps *PacketHandlers) InsertBefore(index int, packetHandlers ...PacketHandlerI) {
	if index <= 0 || index > len(*ps)-1 { // 找不到模板包位置，或者找到第一个，直接插入开头
		tmp := *ps
		*ps = make(PacketHandlers, 0)
		ps.Append(packetHandlers...)
		ps.Append(tmp...)
		return
	}
	before, after := (*ps)[0:index], (*ps)[index:]
	*ps = make(PacketHandlers, 0) // 此处必须重新申请，否则操作会覆盖原有地址
	ps.Append(before...)
	ps.Append(packetHandlers...)
	ps.Append(after...)

}

func (ps *PacketHandlers) InsertAfter(index int, packetHandlers ...PacketHandlerI) {
	if index < 0 || index+1 >= len(*ps) { // 找不到模板包位置,或者目标本就是最后一个，直接在结尾追加
		ps.Append(packetHandlers...)
		return
	}
	before, after := (*ps)[0:index+1], (*ps)[index+1:]
	*ps = make(PacketHandlers, 0) // 此处必须重新申请，否则操作会覆盖原有地址
	ps.Append(before...)
	ps.Append(packetHandlers...)
	ps.Append(after...)
}

func (ps *PacketHandlers) Delete(index int) {
	if index < 0 || len(*ps)-1 < index { // 越界不操作
		return
	}
	if index == len(*ps)-1 { // 需要删除的，在最后一个，直接截断
		*ps = (*ps)[:index]
		return
	}
	before, after := (*ps)[0:index], (*ps)[index+1:]
	*ps = make(PacketHandlers, 0) // 此处必须重新申请，否则操作会覆盖原有地址
	ps.Append(before...)
	ps.Append(after...)
}

func (ps *PacketHandlers) Replace(index int, packetHandlers ...PacketHandlerI) {
	if index < 0 || len(*ps)-1 < index { // 找不到模板包位置,不删除
		return
	}
	before, after := (*ps)[0:index], (*ps)[index+1:]
	*ps = make(PacketHandlers, 0) // 此处必须重新申请，否则操作会覆盖原有地址
	ps.Append(before...)
	ps.Append(packetHandlers...)
	ps.Append(after...)
}

func (ps *PacketHandlers) Index(name string) (indexs []int) {
	indexs = make([]int, 0)
	for i, packet := range *ps {
		if packet.Name() == name {
			indexs = append(indexs, i)
		}
	}
	return indexs
}
func (ps *PacketHandlers) IndexFirst(name string) (index int) {
	for i, packet := range *ps {
		if packet.Name() == name {
			return i
		}
	}
	return -1
}

func (ps *PacketHandlers) IndexLast(name string) (index int) {
	for i := len(*ps) - 1; i < 0; i++ {
		if (*ps)[i].Name() == name {
			return i
		}
	}
	return -1
}
