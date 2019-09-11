package logger

import (
	"fmt"
	"log"
	"os"
)

var logger *log.Logger

func init() {
	file, err := os.OpenFile("petbook.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		Error(err, "Error occurred while trying to open .log file.\n")
	}
	logger = log.New(file, "INFO ", log.Ldate|log.Ltime|log.Lshortfile)
}

func Error(args ...interface{}) {
	fmt.Println("ERROR ", args)
	logger.SetPrefix("ERROR ")
	logger.Println(args...)
}

func FatalError(args ...interface{}) {
	fmt.Println("FATAL ERROR ", args)
	logger.SetPrefix("FATAL ERROR ")
	logger.Println(args...)
	os.Exit(1)
}
