package log

import (
	"fmt"
	"io"
	"log"
	"os"
)

const (
	debugLevel = iota
	infoLevel
	warnLevel
	errLevel
	fatalLevel
)

var logger *Logger

type Logger struct {
	level uint
	err   *log.Logger
	debug *log.Logger
	info  *log.Logger
	warn  *log.Logger
	fatal *log.Logger
}

func init() {
	log.SetFlags(log.Lshortfile | log.Ldate | log.Ltime)
	logger = &Logger{
		level: infoLevel,
		debug: log.New(os.Stdout, "DEBU ", log.Lshortfile|log.Ldate|log.Ltime),
		info:  log.New(os.Stdout, "INFO ", log.Lshortfile|log.Ldate|log.Ltime),
		err:   log.New(os.Stdout, "ERRO ", log.Lshortfile|log.Ldate|log.Ltime),
		warn:  log.New(os.Stdout, "WARN ", log.Lshortfile|log.Ldate|log.Ltime),
		fatal: log.New(os.Stdout, "FATA", log.Lshortfile|log.Ldate|log.Ltime),
	}
}

func SetLevel(level uint) {
	logger.level = level
}

func SetOutput(w io.Writer) {
	logger = &Logger{
		level: logger.level,
		err:   log.New(w, "ERRO ", log.Lshortfile|log.Ldate|log.Ltime),
		debug: log.New(w, "DEBU ", log.Lshortfile|log.Ldate|log.Ltime),
		info:  log.New(w, "INFO ", log.Lshortfile|log.Ldate|log.Ltime),
		warn:  log.New(w, "WARN ", log.Lshortfile|log.Ldate|log.Ltime),
		fatal: log.New(w, "FATA", log.Lshortfile|log.Ldate|log.Ltime),
	}
}

func Debug(msg ...interface{}) {
	if logger.level <= debugLevel {
		logger.debug.Output(2, fmt.Sprint(msg...))
	}
}

func Debugf(format string, v ...interface{}) {
	if logger.level <= debugLevel {
		logger.debug.Output(2, fmt.Sprintf(format, v...))
	}
}

func Info(msg ...interface{}) {
	if logger.level <= infoLevel {
		logger.info.Output(2, fmt.Sprint(msg...))
	}
}

func Infof(format string, v ...interface{}) {
	if logger.level <= infoLevel {
		logger.info.Output(2, fmt.Sprintf(format, v...))
	}
}

func Print(msg ...interface{}) {
	if logger.level <= infoLevel {
		logger.info.Output(2, fmt.Sprint(msg...))
	}
}

func Printf(format string, v ...interface{}) {
	if logger.level <= infoLevel {
		logger.info.Output(2, fmt.Sprintf(format, v...))
	}
}

func Println(msg ...interface{}) {
	if logger.level <= infoLevel {
		logger.info.Output(2, fmt.Sprintln(msg...))
	}
}

func Warn(msg ...interface{}) {
	if logger.level <= warnLevel {
		logger.warn.Output(2, fmt.Sprint(msg...))
	}
}

func Warnf(format string, v ...interface{}) {
	if logger.level <= warnLevel {
		logger.warn.Output(2, fmt.Sprintf(format, v...))
	}
}

func Error(msg ...interface{}) {
	if logger.level <= errLevel {
		logger.err.Output(2, fmt.Sprint(msg...))
	}
}

func Errorf(format string, v ...interface{}) {
	if logger.level <= errLevel {
		logger.err.Output(2, fmt.Sprintf(format, v...))
	}
}

func Fatal(msg ...interface{}) {
	if logger.level <= fatalLevel {
		logger.warn.Output(2, fmt.Sprint(msg...))
	}
	os.Exit(1)
}

func Fatalf(format string, v ...interface{}) {
	if logger.level <= fatalLevel {
		logger.warn.Output(2, fmt.Sprintf(format, v...))
	}
	os.Exit(1)
}
