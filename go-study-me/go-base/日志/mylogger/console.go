package mylogger

import (
	"os"
	"fmt"
)

// ConsoleHandle ...
type ConsoleHandle struct {
	*BaseLogger
}

// NewConsoleLog ...
func NewConsoleLog(level interface{}) *ConsoleHandle {
	lv, ok := parseLogLevel(level)
	if ok {
		panic("make logMgr error...")
	}
	return &ConsoleHandle{BaseLogger: &BaseLogger{lv}}
}

func (c *ConsoleHandle) enable(level LogLevel) bool {
	return level >= c.level
}

func (c *ConsoleHandle) Info(msg string) {
	c.emit(INFO, msg)
}

func (c *ConsoleHandle) Debug(msg string) {
	c.emit(DEBUG, msg)
}

func (c *ConsoleHandle) Error(msg string) {
	c.emit(ERROR, msg)
}

func (c *ConsoleHandle) Warnning(msg string) {
	c.emit(WARNNING, msg)
}

func (c *ConsoleHandle) Fatal(msg string) {
	c.emit(FATAL, msg)
}

func (c *ConsoleHandle) emit(level LogLevel, msg string){
	if c.enable(level) {    
		fmt.Fprintln(os.Stdout, format(level, msg))
	}
}
