package runit

import "fmt"

var (
	// LogLevel is the level at which to start logging.
	// 0 debug
	// 1 info
	// 2 error
	LogLevel = InfoLog
)

const (
	// DebugLog 0
	DebugLog = iota
	// InfoLog 1
	InfoLog
	// ErrorLog 2
	ErrorLog
)

func pdebug(args ...interface{}) {
	if LogLevel <= DebugLog {
		args = append([]interface{}{"[debu]"}, args...)
		fmt.Println(args...)
	}
}

func pdebugf(format string, args ...interface{}) {
	if LogLevel <= DebugLog {
		fmt.Printf("[debu] "+format+"\n", args...)
	}
}

func pinfo(args ...interface{}) {
	if LogLevel <= InfoLog {
		args = append([]interface{}{"[info]"}, args...)
		fmt.Println(args...)
	}
}

func pinfof(format string, args ...interface{}) {
	if LogLevel <= InfoLog {
		fmt.Printf("[info] "+format+"\n", args...)
	}
}

func perror(args ...interface{}) {
	if LogLevel <= ErrorLog {
		args = append([]interface{}{"[erro]"}, args...)
		fmt.Println(args...)
	}
}

func perrorf(format string, args ...interface{}) {
	if LogLevel <= ErrorLog {
		fmt.Printf("[erro] "+format+"\n", args...)
	}
}
