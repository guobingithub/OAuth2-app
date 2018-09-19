package logger

import (
	"runtime"
	"fmt"
	"os"
	"strings"
	"strconv"
	"time"
	"log"
)

const (
	LEV_DEBUG = iota
	LEV_INFO
	LEV_WARN
	LEV_ERROR
	LEV_PANIC
	LEV_FATAL
	TIMFMT    = "2006-01-02 15:04:05"
)

// 2006-01-02 15:04:05.999999
// 日志格式如下
// [2017/07/24 13:49:33] [info] (${packageName} ${fileName} ${methodName}:${lineNum}) inputParams...


// 默认为info
var levelStr = "info"
// 有6种级别
// 1. info
// 2. warn
// 3. debug
// 4. error
// 5. panic 打印完日志后触发系统panic
// 6. fatal 打印完日志后触发os.Exit()
//
// 初始化设置前缀和日志格式
func init() {
	log.SetFlags(log.LUTC)
}

func SetLevel(level string) {
	levelStr = level
}

func getLevel() int32 {
	levelStr := strings.ToLower(levelStr)
	switch levelStr {
	case "debug":
		return LEV_DEBUG
	case "info":
		return LEV_INFO
	case "warn":
		return LEV_WARN
	case "error":
		return LEV_ERROR
	case "panic":
		return LEV_PANIC
	case "fatal":
		return LEV_FATAL
	default:
		return LEV_WARN
	}
}

//获取当前协程的
func Goid() int {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println("panic recover:panic info:%v", err)
		}
	}()

	var buf [64]byte
	n := runtime.Stack(buf[:], false)
	idField := strings.Fields(strings.TrimPrefix(string(buf[:n]), "goroutine "))[0]
	id, err := strconv.Atoi(idField)
	if err != nil {
		panic(fmt.Sprintf("cannot get goroutine id: %v", err))
	}
	return id
}

// 获取调用文件和行号
func caller() string {
	pc, file, line, ok := runtime.Caller(2)
	if !ok {
		line = 0
	}
	pathArgs := strings.Split(file, "/")
	filename := pathArgs[len(pathArgs)-1]
	packAndFuncName := runtime.FuncForPC(pc).Name()
	packFunc := strings.Split(packAndFuncName, ".")
	funcName := packFunc[len(packFunc)-1]
	packName := strings.Join(packFunc[:len(packFunc)-1], ".")
	return fmt.Sprintf("(%s %s %s:%d)", packName, filename, funcName, line)
}

func lPush(level string, caller string, args [] interface{}) []interface{} {
	tempArray := []interface{}{getNow(), level, goIdFmt(), caller}
	tempArray = append(tempArray, args...)
	return tempArray
}

func getNow() string {
	return fmt.Sprint("[", time.Now().Format(TIMFMT), "]")
}

func goIdFmt() string {
	return fmt.Sprintf("[GOID:%d]", Goid())
}

func Debug(v ...interface{}) {
	if getLevel() <= LEV_DEBUG {
		fmt.Println(lPush("[DEBUG]", caller(), v)...)
	}
}

func Info(v ...interface{}) {
	if getLevel() <= LEV_INFO {
		fmt.Println(lPush("[INFO ]", caller(), v)...)
	}
}

func Warn(v ...interface{}) {
	if getLevel() <= LEV_WARN {
		fmt.Println(lPush("[WARN ]", caller(), v)...)
	}
}

func Error(v ...interface{}) {
	if getLevel() <= LEV_ERROR {
		log.Println(lPush("[EROR]", caller(), v)...)
	}
}

func Panic(v ...interface{}) {
	s := caller()
	//os.Exit(1)
	if getLevel() <= LEV_PANIC {
		log.Println(lPush("[PINC_]", s, v)...)
	}
	panic(s)
}

func Fatal(v ...interface{}) {
	if getLevel() <= LEV_FATAL {
		log.Println(lPush("[FATAL]", caller(), v)...)
	}
	os.Exit(1)
}
