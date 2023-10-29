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

type HandlerLog struct {
	Input  []byte `json:"input"`
	Output []byte `json:"output"`
	Err    error
}

type StreamLog struct {
	HandlerLogs []HandlerLog `json:"handlerLogs"`
	logchan.EmptyLogInfo
}

func (w StreamLog) GetName() (logName logchan.LogName) {
	return StreamLogName
}

func (w StreamLog) Error() (err error) {
	for _, handlerLog := range w.HandlerLogs {
		if handlerLog.Err != nil {
			return err
		}
	}
	return nil
}

func DefaultPrintStreamLog(logInfo logchan.LogInforInterface, typeName logchan.LogName, err error) {
	if typeName != StreamLogName {
		return
	}
	streamLog, ok := logInfo.(*StreamLog)
	if !ok {
		return
	}
	processSessionID := logchan.GetSessionID(logInfo)
	if err != nil {
		fmt.Fprintf(logchan.LogWriter, "processSessionID:%s|loginInfo:%s|error:%s\n", processSessionID, streamLog.GetName(), err.Error())
		return
	}

	for i, handlerLog := range streamLog.HandlerLogs {
		fmt.Fprintf(logchan.LogWriter,
			"processSessionID:%s|serialNumber:%d|input:%s|output:%s|err:%s\n",
			processSessionID,
			i,
			string(handlerLog.Input),
			string(handlerLog.Output),
			handlerLog.Err.Error(),
		)
	}
}
