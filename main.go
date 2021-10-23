package logs

import (
	"context"
	"os"
	"time"

	"github.com/sirupsen/logrus"
)

const LOG_TO_FILE = "logToFile"
const LOG_TO_STD_OUT = "logToStdOut"

//Лог файл
var LogToFile *logrus.Logger

//Лог вывод в std output
var LogToStdOut *logrus.Logger

func cyclLogFile(fileLogOpenedChannel chan bool) {
	//дата открытого лог файла.
	var dateLogFile string

	for {
		time.Sleep(1 * time.Second)
		timeNow := time.Now()
		//дата в данный момент.
		dateNow := timeNow.Format("2006-01-02")
		var pathFileLog = "logs/" + dateNow + ".log"
		if dateNow != dateLogFile && dateLogFile != "" {
			openFileLog(LogToFile, pathFileLog, fileLogOpenedChannel)
		}
		if dateLogFile == "" {
			openFileLog(LogToFile, pathFileLog, fileLogOpenedChannel)
		}
		dateLogFile = dateNow
	}
}

//открыть файл для логирования.
func openFileLog(logToFile *logrus.Logger, pathFileLog string, fileLogOpenedChannel chan bool) {
	LogToFile = logrus.New()
	fileLog, err := os.OpenFile(pathFileLog, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0600)
	if err != nil {
		LogToFile.Fatalln(err)
	}
	LogToFile.SetFormatter(&logrus.JSONFormatter{})
	LogToFile.SetOutput(fileLog)
	LogToFile.Infoln("test")
	fileLogOpenedChannel <- true
}

func setLogStdOut() {
	LogToStdOut = logrus.New()
	LogToStdOut.SetOutput(os.Stdout)
}

func New() (*logrus.Logger, *logrus.Logger) {
	fileLogOpenedChannel := make(chan bool)
	go cyclLogFile(fileLogOpenedChannel)
	select {
	case isFileOpened := <-fileLogOpenedChannel:
		if isFileOpened {
			LogToFile.Infoln("File log is open.")
		}
	}

	setLogStdOut()
	return LogToFile, LogToStdOut
}

//получить лог из контекста
func GetLogFromCtx(ctx context.Context, key string) *logrus.Logger {
	switch ctx.Value(key).(type) {
	case *logrus.Logger:
		switch key {
		case LOG_TO_FILE:
			fallthrough
		case LOG_TO_STD_OUT:
			value := ctx.Value(key).(*logrus.Logger)
			return value
		default:
			return nil
		}
	default:
		return nil
	}

}
