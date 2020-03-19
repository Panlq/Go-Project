package mylogger

import (
	"fmt"
	"os"
	"path"
	"time"
)


type logS struct {
	level   LogLevel
	msg string
	funcName string
	fileName string
	timestamp string
	lineno int
}
 
// FileHandleWithChan ...
type FileHandleWithChan struct {
	level LogLevel
	filePath string
	fileName string
	maxFileSize int64
	fileObj *os.File
	errFileObj *os.File
	logChan chan *logS
}
      
// NewFileAsyncHandle ...
func NewFileAsyncHandle(level interface{}, fp, fn string, maxSize int64) *FileHandleWithChan {
	lv, ok := parseLogLevel(level)
	if ok {
		panic("make FileHandleWithChan error...")
	}
	f := &FileHandleWithChan{
		level: lv,
		filePath: fp,
		fileName: fn,
		maxFileSize: maxSize,
		logChan: make(chan *logS, 50000),
	}
	err := f.initFile() // 按照文件路径和文件名将文件打开
	if err != nil {
		panic(err)
	}
	return f
} 

func (f *FileHandleWithChan) initFile() (error) {
	path := path.Join(f.filePath, f.fileName)
	fObj, err := os.OpenFile(path, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Printf("open log file failed, err:%v\n", err)
		return err
	}

	errFObj, err := os.OpenFile(path + ".err", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Printf("open log file failed, err:%v\n", err)
		return err
	}

	// 日志文件已经打开
	f.fileObj = fObj
	f.errFileObj = errFObj
	// 开启后台goroutine写日志
	// for i := 0; i < 5; i++ {
	// 	go f.writeLogBack()
	// }
	go f.writeLogBack()
	return nil

}

func (f *FileHandleWithChan) enable(level LogLevel) bool {
	return level >= f.level
}

func (f *FileHandleWithChan) Info(msg string) {
	f.emit(INFO, msg)
}

func (f *FileHandleWithChan) Debug(msg string) {
	f.emit(DEBUG, msg)
}

func (f *FileHandleWithChan) Error(msg string) {
	f.emit(ERROR, msg)
}

func (f *FileHandleWithChan) Warnning(msg string) {
	f.emit(WARNNING, msg)
}

func (f *FileHandleWithChan) Fatal(msg string) {
	f.emit(FATAL, msg)
}

func (f *FileHandleWithChan) rotateFile(iow *os.File) (*os.File, bool) {
	// 检查文件大小
	fileInfo, err := iow.Stat()
	if err != nil {
		fmt.Printf("get file info failed, err %v\n", err)
		return iow, false
	}
	if fileInfo.Size() < f.maxFileSize {
		return iow, false
	}

	// rename  xx.log --> xx.log.bactimestamp

	nowStr := time.Now().Format("20060102150405000")
	logName := path.Join(f.filePath, fileInfo.Name())
	newLogName := fmt.Sprintf("%s.bak%s", logName, nowStr)
	// 1. 关闭当前日志文件
	iow.Close()
	// 2. 备份一下
	os.Rename(logName, newLogName)
	
	// 3. 打开新文件 并进行重新赋值
	newFObj, err := os.OpenFile(logName, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Printf("open log file failed, err:%v\n", err)
		return iow, false
	}
	return newFObj, true
}

func (f *FileHandleWithChan) writeLogBack() {
	// 后台开线程无限循环取值
	for {
		newFObj, ok := f.rotateFile(f.fileObj)
		if ok {
			f.fileObj = newFObj
		}
	
		select {
		case logTmp := <- f.logChan:
			// get log msg from channel
			logmsg := fmt.Sprintf("[%s]-[%s] [funcName:%s]-[lineno:%d]: %s", 
							logTmp.timestamp, StrLevel[logTmp.level], 
							logTmp.funcName, logTmp.lineno, logTmp.msg)
	
			fmt.Fprintln(f.fileObj, logmsg)
	
			if logTmp.level >= ERROR {
				newFObj, ok := f.rotateFile(f.errFileObj)
				if ok {
					f.errFileObj = newFObj
				}
				errmsg := fmt.Sprintf("[%s]-[%s] [fileName:%s]-[funcName:%s]-[lineno:%d]: %s", 
									logTmp.timestamp, StrLevel[logTmp.level], logTmp.fileName, 
									logTmp.funcName, logTmp.lineno, logTmp.msg)
				fmt.Fprintln(f.errFileObj, errmsg)
			}
		default:
			// 首次会走这里
			time.Sleep(time.Millisecond * 500)
		}
	}
}

func (f *FileHandleWithChan) emit(lev LogLevel, msg string){
	now := time.Now()
	funcName, fileName, lineno := getStackInfo(3)
	if f.enable(lev) {
		logTmp := &logS{
			level: lev,
			msg: msg,
			funcName: funcName,
			fileName: fileName,
			timestamp: now.Format("2006-01-02 15:04:05"),
			lineno: lineno,
		}
		select {
		case f.logChan <- logTmp:
		default:
			fmt.Println("zhixingl zheli ")
			// 通道满了就丢掉保证不出现阻塞, 但是一般是不会出现该情况，因为有另外一个线程从通道中取值
		}
	}
}

// Close ...
func (f *FileHandleWithChan) Close() {
	f.fileObj.Close()
	f.errFileObj.Close()
}