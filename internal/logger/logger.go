package logger

import (
	"encoding/json"
	"fmt"
	"log"
	"reflect"
	"runtime"
	"strings"
	"time"
)

var LogLevel = map[string]int8{
	"TRACE": 0,
	"DEBUG": 1,
	"INFO":  2,
	"WARN":  3,
	"ERROR": 4,
}

type Logger interface {
	Error(message string, args ...interface{})
	Errorf(message string, args ...interface{})
	Warn(message string, args ...interface{})
	Warnf(message string, args ...interface{})
	Info(message string, args ...interface{})
	Infof(message string, args ...interface{})
	Debug(message string, args ...interface{})
	Debugf(message string, args ...interface{})
	Trace(message string, args ...interface{})
	SetTraceId(traceId string)
	PublishSumoLogs()
	Refresh(logLevel string, pushToSumo bool)
}

const (
	TRACE         = "TRACE"
	DEBUG         = "DEBUG"
	INFO          = "INFO"
	WARN          = "WARN"
	ERROR         = "ERROR"
	logFormat     = "%s [%s] [%s] (%s) %s\n"
	logInfoFormat = "%s [%s] [%s] %s\n"
)

func Trace(traceId string, message string, args interface{}) {
	logIt(TRACE, traceId, message, args)
}

func Debug(traceId string, message string, args interface{}) {
	logIt(DEBUG, traceId, message, args)
}

func Info(traceId string, message string, args interface{}) {
	logIt(INFO, traceId, message, args)
}

func Warn(traceId string, message string, args interface{}) {
	logIt(WARN, traceId, message, args)
}

func Error(traceId string, message string, args interface{}) {
	logIt(ERROR, traceId, message, args)
}

func logItf(level string, requestID string, format string, a ...interface{}) {
	msg := fmt.Sprintf(format, a[:len(a)-1]...)
	logIt(level, requestID, msg, a[len(a)-1])
}

func logIt(level string, requestID string, message string, args ...interface{}) {
	src := src()
	now := now()

	var argList = make([]interface{}, 0)

	argList = append(argList, message)
	if args != nil && len(args) > 0 {
		argList = append(argList, jsonIt(args[0]))
	}

	var msg = fmt.Sprint(argList...)

	if level == INFO {
		log.Printf(logInfoFormat, now, level, requestID, msg)
	} else {
		log.Printf(logFormat, now, level, requestID, src, msg)
	}
}

func src() string {
	// Determine caller func
	pc, _, lineno, ok := runtime.Caller(4)
	src := ""
	if ok {
		slice := strings.Split(runtime.FuncForPC(pc).Name(), "/")
		src = slice[len(slice)-1]
		src = fmt.Sprintf("%s:%d", src, lineno)
	}
	return src
}

func now() string {
	return time.Now().Format("2006-01-02T15:04:05-0700")
}

func jsonIt(a interface{}) string {
	dataType := indirectType(reflect.TypeOf(a))
	switch dataType.Kind() {
	case reflect.String, reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Bool,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr,
		reflect.Float32, reflect.Float64, reflect.Complex64, reflect.Complex128:
		return fmt.Sprintf("%v", a)
	}
	data, _ := json.Marshal(a)
	return string(data)
}

func indirectType(reflectType reflect.Type) reflect.Type {
	for isPtrType(reflectType) || reflectType.Kind() == reflect.Slice {
		reflectType = reflectType.Elem()
	}
	return reflectType
}

func isPtrType(t reflect.Type) bool {
	return t.Kind() == reflect.Ptr
}
