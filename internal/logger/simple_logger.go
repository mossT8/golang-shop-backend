package logger

import (
	"fmt"
	"time"
)

type SimpleLogger struct {
	traceId     string
	logLevel    string
	logMessages []map[string]interface{}
	pushToSumo  bool
}

func NewSimpleLogger(logLevel string, pushToSumo bool) *SimpleLogger {
	return &SimpleLogger{
		traceId:    "ROOT",
		logLevel:   logLevel,
		pushToSumo: pushToSumo,
	}
}

func (sl *SimpleLogger) SetTraceId(traceId string) {
	sl.traceId = traceId
}

func (sl *SimpleLogger) Refresh(logLevel string, pushToSumo bool) {
	sl.logLevel = logLevel
	sl.pushToSumo = pushToSumo
	sl.logMessages = nil
}

func (sl *SimpleLogger) Error(message string, args ...interface{}) {
	sl.log(ERROR, message, args...)
}

func (sl *SimpleLogger) Errorf(message string, a ...any) {
	sl.log(ERROR, fmt.Sprintf(message, a...))
}

func (sl *SimpleLogger) Warn(message string, args ...interface{}) {
	sl.log(WARN, message, args...)
}

func (sl *SimpleLogger) Warnf(message string, a ...any) {
	sl.log(WARN, fmt.Sprintf(message, a...))
}

func (sl *SimpleLogger) Info(message string, args ...interface{}) {
	sl.log(INFO, message, args...)
}

func (sl *SimpleLogger) Infof(message string, a ...any) {
	sl.log(INFO, fmt.Sprintf(message, a...))
}

func (sl *SimpleLogger) Debug(message string, args ...any) {
	sl.log(DEBUG, message, args...)
}

func (sl *SimpleLogger) Debugf(message string, a ...any) {
	sl.log(DEBUG, fmt.Sprintf(message, a...))
}

func (sl *SimpleLogger) Trace(message string, args ...interface{}) {
	sl.log(TRACE, message, args...)
}

func (sl *SimpleLogger) Log(traceId string, level string, message string, args ...interface{}) {
	sl.SetRequestId(traceId)
	sl.log(level, message, args...)
}

func (sl *SimpleLogger) LogLevel() string {
	return sl.logLevel
}

func (sl *SimpleLogger) LogMessages() []map[string]interface{} {
	return sl.logMessages
}

func (sl *SimpleLogger) RequestId() string {
	return sl.traceId
}

func (sl *SimpleLogger) log(level string, message string, args ...interface{}) {
	if LogLevel[level] >= LogLevel[sl.logLevel] {
		logIt(level, sl.traceId, message, args...)
		if sl.pushToSumo {
			sl.appendMessageToList(level, message, args...)
		}
	}
}

func (sl *SimpleLogger) appendMessageToList(level string, message string, args ...interface{}) {

	timestamp := time.Now().Format("2006-01-02 15:04:05.000")

	var payload = make(map[string]interface{})
	if LogLevel[level] == 2 || LogLevel[level] == 3 { //Dont add source for INFO & WARN log levels
		payload[level] = fmt.Sprintf("[%s] [%s] - %s", timestamp, sl.traceId, message)
	} else {
		payload[level] = fmt.Sprintf("[%s] [%s] (%s) - %s", timestamp, sl.traceId, src(), message)
	}

	if args != nil && len(args) > 0 {
		payload["msg-args"] = args[0]
	}
	sl.logMessages = append(sl.logMessages, payload)
}

func (sl *SimpleLogger) PublishSumoLogs() {

	// if sl.pushToSumo && sl.sqsProducer != nil {

	// 	batchSize, err := strconv.Atoi(sl.config.BatchSize)
	// 	if err != nil {
	// 		Warn(sl.traceId, "Invalid sumo config batch-size value: ", err.Error())
	// 		batchSize = 15
	// 	}

	// 	var batchCounter = 0
	// 	var logBatch = make([]types.SendMessageBatchRequestEntry, 0)
	// 	for i := 0; i < len(sl.logMessages); i++ {

	// 		messageId := uuid.NewString()
	// 		messageBody := jsonIt(map[string]interface{}{
	// 			"host":        sl.config.OnboardHost,
	// 			"category":    sl.config.Category,
	// 			"sumoPayload": sl.logMessages[i],
	// 		})

	// 		logBatch = append(logBatch, types.SendMessageBatchRequestEntry{
	// 			Id:          &messageId,
	// 			MessageBody: &messageBody,
	// 		})
	// 		batchCounter++

	// 		if batchCounter == batchSize || i == len(sl.logMessages)-1 {

	// 			Debug(sl.traceId, "Publishing log batch to sumo: ", logBatch)
	// 			err := sl.sqsProducer.SendMessageBatch(&logBatch)
	// 			if err != nil {
	// 				Error(sl.traceId, "Failed to publish sumo log batch: ", err.Error())
	// 			}
	// 			logBatch = make([]types.SendMessageBatchRequestEntry, 0)
	// 			batchCounter = 0
	// 		}
	// 	}
	// }
	sl.logMessages = nil
}

func (sl *SimpleLogger) SetRequestId(traceId string) {
	sl.traceId = traceId
}

func (sl *SimpleLogger) setTimeZoneLocation(timeZoneLocation string) {

	var location = "Africa/Johannesburg"
	if timeZoneLocation != "" {
		location = timeZoneLocation
	}

	loc, err := time.LoadLocation(location)
	if err != nil {
		Warn(sl.traceId, fmt.Sprintf("Failed to load timezone location '%s': ", location), err.Error())
		return
	}
	time.Local = loc
}
