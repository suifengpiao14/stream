package stream

import (
	"fmt"

	"github.com/suifengpiao14/logchan/v2"
)

const (
	StreamLogName LogName = "streamLog"
)

type LogName string

func (l LogName) String() string {
	return string(l)
}

const (
	HandlerLog_Type_SetContext = "setContext"
	HandlerLog_Type_Before     = "before"
	HandlerLog_Type_After      = "after"
)

type HandlerLog struct {
	PackName string
	Order    int
	Type     string // before |after
	Input    []byte `json:"input"`
	Output   []byte `json:"output"`
	Err      error
}

type StreamLog struct {
	HandlerLogs []HandlerLog `json:"handlerLogs"`
	logchan.EmptyLogInfo
}

func (w StreamLog) GetName() (logName logchan.LogName) {
	return StreamLogName
}

func (w StreamLog) Error() (err error) {
	for _, packHandlerLog := range w.HandlerLogs {
		if packHandlerLog.Err != nil {
			return packHandlerLog.Err
		}
	}
	return nil
}

const ColorRed = "\033[0;31m"
const ColorNone = "\033[0m"

func DefaultPrintStreamLog(logInfo logchan.LogInforInterface, typeName logchan.LogName, err error) {
	if typeName != StreamLogName {
		return
	}
	streamLog, ok := logInfo.(*StreamLog)
	if !ok {
		return
	}
	processSessionID := logchan.GetSessionID(logInfo)
	// 把没有出错的步骤日志输出，方便定位问题
	// if err != nil {
	// 	fmt.Fprintf(logchan.LogWriter, "processSessionID:%s|loginInfo:%s|error:%s\n", processSessionID, streamLog.GetName(), err.Error())
	// 	return
	// }

	for i, handlerLog := range streamLog.HandlerLogs {
		errStr := ""
		if handlerLog.Err != nil {
			errStr = handlerLog.Err.Error()
		}
		fmt.Fprintf(logchan.LogWriter,
			"processSessionID:%s|name:%s|serialNumber:%d|type:%s|input:%s|output:%s|err:%s%s%s\n",
			processSessionID,
			handlerLog.PackName,
			i,
			handlerLog.Type,
			string(handlerLog.Input),
			string(handlerLog.Output),
			ColorRed,
			errStr,
			ColorNone,
		)
	}
}
