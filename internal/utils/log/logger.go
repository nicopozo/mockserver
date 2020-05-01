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

func (theLogger *log) Info(source interface{}, tags map[string]string, message string, args ...interface{}) {
	theLogger.logrus.Infof("%s - %v", theLogger.GetMessage(message, args...), theLogger.getTags(source, tags))
}

func (theLogger *log) Warn(source interface{}, tags map[string]string, message string, args ...interface{}) {
	theLogger.logrus.Warnf("%s - %v", theLogger.GetMessage(message, args...), theLogger.getTags(source, tags))
}

func (theLogger *log) Error(source interface{}, tags map[string]string, err error,
	message string, args ...interface{}) {
	theLogger.logrus.Errorf("%s - error: %s - %v", theLogger.GetMessage(message, args...), err.Error(),
		theLogger.getTags(source, tags))
}

func (theLogger *log) Debug(source interface{}, tags map[string]string, message string, args ...interface{}) {
	theLogger.logrus.Debugf("%s - %v", theLogger.GetMessage(message, args...), theLogger.getTags(source, tags))
}

func (theLogger *log) GetTrackingID() string {
	return theLogger.trackingID
}

func (theLogger *log) GetMessage(message string, args ...interface{}) string {
	if len(args) > 0 {
		return fmt.Sprintf(message, args...)
	}

	return message
}

func newRequestID() string {
	id := ""
	logID, err := uuid.NewV4()

	if err == nil {
		id = logID.String()
	}

	return id
}

func getClass(source interface{}) string {
	t := reflect.TypeOf(source)
	if t != nil {
		return t.String()
	}

	return ""
}

func (theLogger *log) getTags(source interface{}, tags map[string]string) []string {
	var res []string

	i := 0

	if len(tags) == 0 {
		res = make([]string, minTags)
	} else {
		res = make([]string, len(tags)+minTags)
		for key, value := range tags {
			res[i] = fmt.Sprintf("%s:%v", key, value)
			i++
		}
	}

	res[i] = fmt.Sprintf("TRACKING_ID:%v", theLogger.trackingID)
	res[i+1] = fmt.Sprintf("Class:%v", getClass(source))

	return res
}
