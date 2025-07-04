package logs

import (
	"os"
	"path/filepath"

	"github.com/sirupsen/logrus"
)

var Logger = logrus.New()

func InitLogger() {
	logDir := "logs"
	os.MkdirAll(logDir, os.ModePerm)


	allLog, err := os.OpenFile(filepath.Join(logDir, "all.log"), os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		logrus.Fatalf("No se pudo crear archivo de log: %v", err)
	}

	errorLog, err := os.OpenFile(filepath.Join(logDir, "error.log"), os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		logrus.Fatalf("No se pudo crear archivo de error: %v", err)
	}

	Logger.SetFormatter(&logrus.JSONFormatter{})
	Logger.SetOutput(allLog) 


	Logger.AddHook(&ErrorHook{
		Writer: errorLog,
		LogLevels: []logrus.Level{
			logrus.ErrorLevel,
			logrus.FatalLevel,
			logrus.PanicLevel,
		},
	})
}

type ErrorHook struct {
	Writer    *os.File
	LogLevels []logrus.Level
}

func (hook *ErrorHook) Fire(entry *logrus.Entry) error {
	line, err := entry.Bytes()
	if err != nil {
		return err
	}
	_, err = hook.Writer.Write(line)
	return err
}

func (hook *ErrorHook) Levels() []logrus.Level {
	return hook.LogLevels
}
