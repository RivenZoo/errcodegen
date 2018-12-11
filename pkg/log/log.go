package log

import (
	"log"
)

func Info(fmtStr string, args ...interface{}) {
	log.Printf("[INFO] "+fmtStr+"\n", args...)
}

func Debug(fmtStr string, args ...interface{}) {
	log.Printf("[DEBUG] "+fmtStr+"\n", args...)
}

func Error(fmtStr string, args ...interface{}) {
	log.Printf("[ERROR] "+fmtStr+"\n", args...)
}
