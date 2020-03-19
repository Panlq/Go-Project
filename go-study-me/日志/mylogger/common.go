package mylogger

import (
	"fmt"
	"time"
	"path"
	"errors"
	"runtime"
	"strings"
)

//LogLevel ...
type LogLevel uint16

const (
	UNKNOW          = LogLevel(0)
	INFO   LogLevel = iota * 10
	DEBUG
	ERROR
	WARNNING
	FATAL
)

// StrLevel ...
var StrLevel = map[LogLevel]string{
	INFO: "INFO",
	DEBUG: "DEBUG",
	WARNNING: "WARNNING",
	ERROR: "ERROR",
	FATAL: "FATAL",
}

func parseLogLevel(x interface{}) (LogLevel, bool) {
	var err bool = false
	var level LogLevel = UNKNOW
	if v, ok := x.(string); ok {
		s := strings.ToLower(v)
		switch s {
		case "debug":
			level = DEBUG
		case "info":
			level = INFO
		case "error":
			level = ERROR
		case "warn":
			level = WARNNING
		case "fatal":
			level = FATAL
		default:
			err = true
		}

	} else if v, ok := x.(LogLevel); ok {
		level = v
	} else {
		err = true
	}

	if err {
		errs := errors.New("无效的日志级别")
		if errs != nil {
			fmt.Println(errs)
		}
	}

	return level, err
}


func format(level LogLevel, msg string) string {
	now := time.Now()
	ts := now.Format("2006-01-02 15:04:05")
	funcName, fileName, lineno := getStackInfo(4)
	if level >= ERROR {
		return fmt.Sprintf("[%s]-[%s] [fileName:%s]-[funcName:%s]-[lineno:%d]: %s", ts, StrLevel[level], fileName, funcName, lineno, msg)
	}
	return fmt.Sprintf("[%s]-[%s] [funcName:%s]-[lineno:%d]: %s", ts, StrLevel[level], funcName, lineno, msg)
}


//get line-no, filename, 
func getStackInfo(skip int) (funcName, fileName string, lineno int){
	pc, file, lineno, ok := runtime.Caller(skip)
	funcName = runtime.FuncForPC(pc).Name()
	fileName = path.Base(file)
	if !ok {
		fmt.Printf("runtime.Caller() fai ")
	}
	return
}

type BaseLogger struct {
	level LogLevel
}
