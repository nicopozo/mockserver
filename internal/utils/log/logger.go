package log

import (
	"fmt"
	"os"
	"reflect"
	"time"

	"github.com/gofrs/uuid"
	"github.com/sirupsen/logrus"
)

const (
	minTags = 2
)

type ILogger interface {
	Info(source interface{}, tags map[string]string, message string, args ...interface{})
	Warn(source interface{}, tags map[string]string, message string, args ...interface{})
	Error(source interface{}, tags map[string]string, err error, message string, args ...interface{})
	Debug(source interface{}, tags map[string]string, message string, args ...interface{})
	GetTrackingID() string
	GetMessage(message string, args ...interface{}) string
}

type log struct {
	trackingID string
	logrus     *logrus.Logger
}

func DefaultLogger() ILogger {
	format := new(logrus.TextFormatter)
	format.TimestampFormat = time.RFC3339
	format.FullTimestamp = true
	iLogger := &log{
		trackingID: newRequestID(),
		logrus: &logrus.Logger{
			Out:       os.Stdout,
			Formatter: format,
			Hooks:     make(logrus.LevelHooks),
			Level:     logrus.DebugLevel,
		},
	}

	return iLogger
}

func NewLogger(trackingID string) ILogger {
	format := new(logrus.TextFormatter)
	format.TimestampFormat = time.RFC3339
	format.FullTimestamp = true
	iLogger := &log{
		trackingID: trackingID,
		logrus: &logrus.Logger{
			Out:       os.Stdout,
			Formatter: format,
			Hooks:     make(logrus.LevelHooks),
			Level:     logrus.DebugLevel,
		},
	}

	return iLogger
}

func (logger *log) Info(source interface{}, tags map[string]string, message string, args ...interface{}) {
	logger.logrus.Infof("%s - %v", logger.GetMessage(message, args...), logger.getTags(source, tags))
}

func (logger *log) Warn(source interface{}, tags map[string]string, message string, args ...interface{}) {
	logger.logrus.Warnf("%s - %v", logger.GetMessage(message, args...), logger.getTags(source, tags))
}

func (logger *log) Error(source interface{}, tags map[string]string, err error, message string, args ...interface{}) {
	logger.logrus.Errorf("%s - error: %s - %v", logger.GetMessage(message, args...), err.Error(),
		logger.getTags(source, tags))
}

func (logger *log) Debug(source interface{}, tags map[string]string, message string, args ...interface{}) {
	logger.logrus.Debugf("%s - %v", logger.GetMessage(message, args...), logger.getTags(source, tags))
}

func (logger *log) GetTrackingID() string {
	return logger.trackingID
}

func (logger *log) GetMessage(message string, args ...interface{}) string {
	if len(args) > 0 {
		return fmt.Sprintf(message, args...)
	}

	return message
}

func newRequestID() string {
	requestID := ""
	logID, err := uuid.NewV4()

	if err == nil {
		requestID = logID.String()
	}

	return requestID
}

func getClass(source interface{}) string {
	if t := reflect.TypeOf(source); t != nil {
		return t.String()
	}

	return ""
}

func (logger *log) getTags(source interface{}, tags map[string]string) []string {
	var res []string

	index := 0

	if len(tags) == 0 {
		res = make([]string, minTags)
	} else {
		res = make([]string, len(tags)+minTags)
		for key, value := range tags {
			res[index] = fmt.Sprintf("%s:%v", key, value)
			index++
		}
	}

	res[index] = fmt.Sprintf("TRACKING_ID:%v", logger.trackingID)
	res[index+1] = fmt.Sprintf("Class:%v", getClass(source))

	return res
}
