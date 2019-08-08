package logger

import (
	"fmt"
	"log"
	"os"
)

var Logger *log.Logger

func init() {
	file, err := os.OpenFile("petbook.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalln("Failed to open log file", err)
	}
	Logger = log.New(file, "INFO ", log.Ldate|log.Ltime|log.Lshortfile)
}

func Error(args ...interface{}) {
	fmt.Println(args...)
	Logger.SetPrefix("ERROR ")
	Logger.Println(args...)
}

func FatalError(args ...interface{}) {
	fmt.Println(args...)
	Logger.SetPrefix("FATAL ERROR ")
	Logger.Println(args...)
	os.Exit(3)
}