package stream

import (
	"context"
	"fmt"

	"github.com/suifengpiao14/funcs"
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

// 定回调函数指针的类型
type HandlerFn func(ctx context.Context, input []byte) (out []byte, err error)
type ErrorHandler func(ctx context.Context, err error) (out []byte)

type SetContextFn func(ctx context.Context, input []byte) (newCtx context.Context, out []byte, err error)

type PackHandler struct {
	Name       string
	SetContext SetContextFn
	Before     HandlerFn
	After      HandlerFn
}

func NewPackHandler(before HandlerFn, after HandlerFn) (p PackHandler) {
	p = NewPackHandlerWithSetContext(nil, before, after)
	return p
}
func NewPackHandlerWithSetContext(setContextFn SetContextFn, before HandlerFn, after HandlerFn) (p PackHandler) {
	p = PackHandler{
		SetContext: setContextFn,
		Name:       fmt.Sprintf("%s-%s", funcs.GetFuncname(before), funcs.GetFuncname(after)),
		Before:     before,
		After:      after,
	}
	return p
}

type PackHandlers []PackHandler

func (ps *PackHandlers) Add(packHandlers ...PackHandler) {
	if *ps == nil {
		*ps = make(PackHandlers, 0)
	}
	*ps = append(*ps, packHandlers...)
}

type StreamI interface {
	Run(ctx context.Context, input []byte) (out []byte, err error)
	AddPack(handlerPacks ...PackHandler) // 包裹更多处理函数
}

// 任务节点结构定义
type Stream struct {
	packHandlers PackHandlers // 处理链条集合
	errorHandler ErrorHandler //错误处理
}

func NewStream(errorHandelr ErrorHandler, packHandlers ...PackHandler) *Stream {
	stream := &Stream{
		packHandlers: packHandlers,
		errorHandler: errorHandelr,
	}
	return stream
}

// AddPack 增加打包
func (s *Stream) AddPack(handlerPacks ...PackHandler) {
	s.packHandlers = append(s.packHandlers, handlerPacks...)
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
	l := len(s.packHandlers)
	streamLog := StreamLog{
		HandlerLogs: make([]HandlerLog, 0),
	}
	defer func() {
		logchan.SendLogInfo(&streamLog)
	}()
	for i := 0; i < l; i++ { // 先执行最后的before，直到最早的before
		pack := s.packHandlers[i]
		if pack.SetContext != nil {
			handlerLog := HandlerLog{
				Input:    data,
				PackName: pack.Name,
				Type:     HandlerLog_Type_SetContext,
			}
			ctx, input, err = pack.SetContext(ctx, input)
			if err != nil {
				return nil, err
			}
			handlerLog.Err = err
			streamLog.HandlerLogs = append(streamLog.HandlerLogs, handlerLog)
		}
		if pack.Before != nil {
			handlerLog := HandlerLog{
				Input:    data,
				PackName: pack.Name,
				Type:     HandlerLog_Type_Before,
			}
			data, err = pack.Before(ctx, data)
			if err != nil {
				return nil, err
			}
			//handlerLog.Output = data
			handlerLog.Err = err
			streamLog.HandlerLogs = append(streamLog.HandlerLogs, handlerLog)

			if err != nil {
				return nil, err
			}
		}
	}

	for i := l - 1; i > -1; i-- { // 先执行最后的after，直到最早的after
		pack := s.packHandlers[i]
		if pack.After != nil {
			handlerLog := HandlerLog{
				Input:    data,
				PackName: pack.Name,
				Type:     HandlerLog_Type_After,
			}
			handlerLog.Input = data
			data, err = pack.After(ctx, data)
			if err != nil {
				return nil, err
			}
			//handlerLog.Output = data
			handlerLog.Err = err
			streamLog.HandlerLogs = append(streamLog.HandlerLogs, handlerLog)
			if err != nil {
				return nil, err
			}
		}
	}
	return data, err
}
