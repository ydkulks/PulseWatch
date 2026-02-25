package logger

import (
	"io"
	"log"
	"os"
)

func Init() *os.File {
	file, err := os.OpenFile("PulseWatch.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatal(err)
	}
	// defer file.Close()

	multiWriter := io.MultiWriter(os.Stderr, file)
	log.SetOutput(multiWriter)
	return file
}

func Close(file *os.File) {
	file.Close()
}
