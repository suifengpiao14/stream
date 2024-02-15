package stream

import (
	"context"

	"github.com/pkg/errors"
	"github.com/suifengpiao14/logchan/v2"
)

/**
流式工作原理：
各个任务都过指针链表的方式组成一个任务链，这个任务链从第一个开始执行，直到最后一个
每一个任务节点执行完毕会将结果带入到下一级任务节点中。
每一个任务是一个Stream节点，每个任务节点都包含首节点和下一个任务节点的指针,
除了首节点，每个节都会设置一个回调函数的指针，用本节点的任务执行,
最后一个节点的nextStream为空,表示任务链结束。
**/

type StreamI interface {
	Run(ctx context.Context, input []byte) (out []byte, err error)
	AddPack(handlerPacks ...PacketHandlerI) // 包裹更多处理函数
}

// 任务节点结构定义
type Stream struct {
	packetHandlers PacketHandlers // 处理链条集合
	errorHandler   ErrorHandler   //错误处理
}

func NewStream(name string, errorHandelr ErrorHandler, packetHandlers ...PacketHandlerI) *Stream {
	stream := &Stream{
		packetHandlers: packetHandlers,
		errorHandler:   errorHandelr,
	}
	return stream
}

// AddPack 增加打包
func (s *Stream) AddPack(handlerPacks ...PacketHandlerI) {
	s.packetHandlers = append(s.packetHandlers, handlerPacks...)
}

func (s *Stream) Run(ctx context.Context, input []byte) (out []byte, err error) {
	out, err = s.run(ctx, input)
	if err != nil && s.errorHandler != nil {
		out = s.errorHandler(ctx, err)
		err = nil
	}

	if err != nil {
		return nil, err
	}
	return out, nil

}
func (s *Stream) run(ctx context.Context, input []byte) (out []byte, err error) {
	data := input
	l := len(s.packetHandlers)
	streamLog := StreamLog{
		HandlerLogs: make([]HandlerLog, 0),
	}
	defer func() {
		streamLog.SetContext(ctx)
		logchan.SendLogInfo(&streamLog)
	}()
	for i := 0; i < l; i++ { // 先执行最后的before，直到最早的before
		pack := s.packetHandlers[i]
		handlerLog := HandlerLog{
			BeforeCtx: ctx,
			Input:     data,
			PackName:  pack.Name(),
			Type:      HandlerLog_Type_Before,
			Serialize: pack.String(),
		}
		ctx, data, err = pack.Before(ctx, data)
		if errors.Is(err, ERROR_EMPTY_FUNC) { // 这个错误标记是空函数，可以不计日志
			err = nil
			continue
		}
		handlerLog.Err = err
		handlerLog.AfterCtx = ctx
		handlerLog.Output = data
		streamLog.HandlerLogs = append(streamLog.HandlerLogs, handlerLog)
		if err != nil {
			return nil, err
		}
	}

	for i := l - 1; i > -1; i-- { // 先执行最后的after，直到最早的after
		pack := s.packetHandlers[i]
		handlerLog := HandlerLog{
			BeforeCtx: ctx,
			Input:     data,
			PackName:  pack.Name(),
			Type:      HandlerLog_Type_After,
			Serialize: pack.String(),
		}
		handlerLog.Input = data
		ctx, data, err = pack.After(ctx, data)
		if errors.Is(err, ERROR_EMPTY_FUNC) { // 这个错误标记是空函数，可以不计日志
			err = nil
			continue
		}
		handlerLog.Err = err
		handlerLog.AfterCtx = ctx
		handlerLog.Output = data
		streamLog.HandlerLogs = append(streamLog.HandlerLogs, handlerLog)
		if err != nil {
			return nil, err
		}
	}
	return data, err
}
