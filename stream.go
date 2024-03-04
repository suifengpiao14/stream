package stream

import (
	"context"
	"reflect"

	"github.com/pkg/errors"
	"github.com/suifengpiao14/packethandler"
)

/**
流式工作原理：
各个任务都过指针链表的方式组成一个任务链，这个任务链从第一个开始执行，直到最后一个
每一个任务节点执行完毕会将结果带入到下一级任务节点中。
每一个任务是一个Stream节点，每个任务节点都包含首节点和下一个任务节点的指针,
除了首节点，每个节都会设置一个回调函数的指针，用本节点的任务执行,
最后一个节点的nextStream为空,表示任务链结束。
**/
type ErrorHandler func(ctx context.Context, err error) (out []byte)
type StreamI interface {
	Run(ctx context.Context, input []byte) (out []byte, err error)
	AddPack(handlerPacks ...packethandler.PacketHandlerI) // 包裹更多处理函数
}

// 任务节点结构定义
type Stream struct {
	Context        context.Context
	Name           string
	packetHandlers packethandler.PacketHandlers // 处理链条集合
	errorHandler   ErrorHandler                 //错误处理
}

func NewStream(name string, errorHandler ErrorHandler, packetHandlers ...packethandler.PacketHandlerI) *Stream {
	stream := &Stream{
		packetHandlers: packetHandlers,
		errorHandler:   errorHandler,
	}
	return stream
}

func (s *Stream) SetContextValue(key any, value any) {
	if s.Context == nil {
		s.Context = context.Background()
	}
	s.Context = context.WithValue(s.Context, key, value)
}

func (s *Stream) GetContextValue(key any, dest any) (err error) {
	if s.Context == nil {
		s.Context = context.Background()
	}
	anyValue := s.Context.Value(key)
	destRv := reflect.Indirect(reflect.ValueOf(dest))
	anyRV := reflect.Indirect(reflect.ValueOf(anyValue))
	destRT := destRv.Type()
	if !anyRV.CanConvert(destRT) {
		err = errors.Errorf("type error want:%s,got:%s", destRT.Name(), anyRV.Type().Name())
		return err
	}
	anyRV = anyRV.Convert(destRT)
	if !destRv.CanSet() {
		err = errors.Errorf("dst can not set")
		return err
	}
	destRv.Set(anyRV)
	return nil
}

// AddPack 增加打包
func (s *Stream) AddPack(handlerPacks ...packethandler.PacketHandlerI) {
	s.packetHandlers = append(s.packetHandlers, handlerPacks...)
}

func (s *Stream) Run(ctx context.Context, input []byte) (out []byte, err error) {
	out, err = s.packetHandlers.Run(ctx, input)
	if err != nil && s.errorHandler != nil {
		out = s.errorHandler(ctx, err)
		err = nil
	}
	if err != nil {
		return nil, err
	}
	return out, nil

}
