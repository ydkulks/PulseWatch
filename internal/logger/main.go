package logger

import (
	"io"
	"log"
	"os"
)

var (
	infoLogger *log.Logger
	warnLogger *log.Logger
	errorLogger *log.Logger
)

func Init() *os.File {
	file, err := os.OpenFile("PulseWatch.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatal(err)
	}
	// defer file.Close()

	multiWriter := io.MultiWriter(os.Stderr, file)
	log.SetOutput(multiWriter)

	log.SetPrefix("[PulseWatch] ")
	flags := log.Ldate | log.Ltime | log.Llongfile

	infoLogger = log.New(multiWriter, "[INFO]	", flags)
	warnLogger = log.New(multiWriter, "[WARN]	", flags)
	errorLogger = log.New(multiWriter, "[ERROR]	", flags)

	return file
}

func Close(file *os.File) {
	_ = file.Close()
}

func Info(v ...any) {
	infoLogger.Println(v...)
}

func Warn(v ...any) {
	warnLogger.Println(v...)
}

func Error(v ...any) {
	errorLogger.Println(v...)
}
