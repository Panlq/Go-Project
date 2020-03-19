package mylogger

import (
	"fmt"
	"os"
	"path"
	"time"
)


// type logs struct {
// 	level   LogLevel
// 	msg string
// }   

// var logChan = make(chan *logs, 100)
 
// FileHandle ...
type FileHandle struct {
	level LogLevel
	filePath string
	fileName string
	maxFileSize int64
	fileObj *os.File
	errFileObj *os.File
}
      
// NewFileHandle ...
func NewFileHandle(level interface{}, fp, fn string, maxSize int64) *FileHandle {
	lv, ok := parseLogLevel(level)
	if ok {
		panic("make FileHandle error...")
	}
	f := &FileHandle{
		level: lv,
		filePath: fp,
		fileName: fn,
		maxFileSize: maxSize,
	}
	err := f.initFile() // 按照文件路径和文件名将文件打开
	if err != nil {
		panic(err)
	}
	// go f.consumelog()
	return f
} 

func (f *FileHandle) initFile() (error) {
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
	return nil

}

func (f *FileHandle) enable(level LogLevel) bool {
	return level >= f.level
}

// with channel
// func (f *FileHandle) Info(msg string) {
// 	appendlog(&logs{INFO, msg})
// }

// func (f *FileHandle) Debug(msg string) {
// 	appendlog(&logs{DEBUG, msg})
// }

// func (f *FileHandle) Error(msg string) {
// 	appendlog(&logs{ERROR, msg})
// }

// func (f *FileHandle) Warnning(msg string) {
// 	appendlog(&logs{WARNNING, msg})
// }

// func (f *FileHandle) Fatal(msg string) {
// 	appendlog(&logs{FATAL, msg})
// }

// func appendlog(l *logs) {
// 	logChan <- l
// }

// func (f *FileHandle) consumelog() {
// 	for l := range logChan{
// 		fmt.Printf("emit log: level:%s msg:%s\n", StrLevel[l.level], l.msg)
// 		f.emit(l.level, l.msg)
// 	}
// }

func (f *FileHandle) Info(msg string) {
	f.emit(INFO, msg)
}

func (f *FileHandle) Debug(msg string) {
	f.emit(DEBUG, msg)
}

func (f *FileHandle) Error(msg string) {
	f.emit(ERROR, msg)
}

func (f *FileHandle) Warnning(msg string) {
	f.emit(WARNNING, msg)
}

func (f *FileHandle) Fatal(msg string) {
	f.emit(FATAL, msg)
}
 
// func (f *FileHandle) checkSize(fobj *os.File) (bool) {
// 	// get the file size
// 	fileInfo, err := fobj.Stat()
// 	if err != nil {
// 		fmt.Printf("get file info failed, err %v\n", err)
// 		return false
// 	}
// 	return fileInfo.Size() >= f.maxFileSize
// }

func (f *FileHandle) rotateFile(iow *os.File) (*os.File, bool) {
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

func (f *FileHandle) emit(level LogLevel, msg string){
	if f.enable(level) {
		newFObj, ok := f.rotateFile(f.fileObj)
		if ok {
			f.fileObj = newFObj
		}
		iow := f.fileObj
		if level >= ERROR {
			newFObj, ok := f.rotateFile(f.errFileObj)
			if ok {
				f.errFileObj = newFObj
			}
			iow = f.errFileObj
		}

		fmt.Fprintln(iow, format(level, msg))
	}
}

func (f *FileHandle) Close() {
	f.fileObj.Close()
	f.errFileObj.Close()
}