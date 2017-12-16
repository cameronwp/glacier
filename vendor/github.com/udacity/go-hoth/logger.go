package hoth

import (
	"fmt"
	"log"
)

// Logger interface makes logging pluggable for all hoth-api based API packages.
// Required functions are the intersection of functions from logrus and go-logging.
// logrus: https://github.com/Sirupsen/logrus/blob/4b6ea7319e214d98c938f12692336f7ca9348d6b/logrus.go
// go-logging: https://github.com/op/go-logging/blob/e6cee1331aceff5b42d8a3fe6c585e30330dcf39/logger.go
type Logger interface {
	Debugf(format string, args ...interface{})
	Infof(format string, args ...interface{})
	Warningf(format string, args ...interface{})
	Errorf(format string, args ...interface{})
	Fatalf(format string, args ...interface{})
	Panicf(format string, args ...interface{})

	Debug(args ...interface{})
	Info(args ...interface{})
	Warning(args ...interface{})
	Error(args ...interface{})
	Fatal(args ...interface{})
	Panic(args ...interface{})
}

// DefaultLogger implements the Logger interface using stdlib log.
type DefaultLogger struct{}

// Debugf outputs to the standard logger with a DEBUG prefix.
func (l *DefaultLogger) Debugf(format string, args ...interface{}) {
	log.Printf("DEBU: "+format, args...)
}

// Infof outputs to the standard logger with a INFO prefix.
func (l *DefaultLogger) Infof(format string, args ...interface{}) {
	log.Printf("INFO: "+format, args...)
}

// Warningf outputs to the standard logger with a WARN prefix.
func (l *DefaultLogger) Warningf(format string, args ...interface{}) {
	log.Printf("WARN: "+format, args...)
}

// Errorf outputs to the standard logger with a ERRO prefix.
func (l *DefaultLogger) Errorf(format string, args ...interface{}) {
	log.Printf("ERRO: "+format, args...)
}

// Fatalf outputs to the standard logger Fatalf() with a FATA prefix.
func (l *DefaultLogger) Fatalf(format string, args ...interface{}) {
	log.Fatalf("FATA: "+format, args...)
}

// Panicf outputs to the standard logger Panicf() with a PANI prefix.
func (l *DefaultLogger) Panicf(format string, args ...interface{}) {
	log.Panicf("PANI: "+format, args...)
}

// Debug outputs to the standard logger with a DEBUG prefix.
func (l *DefaultLogger) Debug(args ...interface{}) {
	log.Print("DEBU: ", fmt.Sprint(args...))
}

// Info outputs to the standard logger with a INFO prefix.
func (l *DefaultLogger) Info(args ...interface{}) {
	log.Print("INFO: ", fmt.Sprint(args...))
}

// Warning outputs to the standard logger with a WARN prefix.
func (l *DefaultLogger) Warning(args ...interface{}) {
	log.Print("WARN: ", fmt.Sprint(args...))
}

// Error outputs to the standard logger with a ERRO prefix.
func (l *DefaultLogger) Error(args ...interface{}) {
	log.Print("ERRO: ", fmt.Sprint(args...))
}

// Fatal outputs to the standard logger Fatal() with a FATA prefix.
func (l *DefaultLogger) Fatal(args ...interface{}) {
	log.Fatal("FATA: ", fmt.Sprint(args...))
}

// Panic outputs to the standard logger Panic() with a PANI prefix.
func (l *DefaultLogger) Panic(args ...interface{}) {
	log.Panic("PANI: ", fmt.Sprint(args...))
}
