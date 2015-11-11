package runit

import "fmt"

var (
	LogLevel = InfoLog
)

const (
	DebugLog = iota
	InfoLog
	ErrorLog
)

func pdebug(args ...interface{}) {
	if LogLevel <= DebugLog {
		args = append([]interface{}{"debu]"}, args...)
		fmt.Println(args...)
	}
}

func pdebugf(format string, args ...interface{}) {
	if LogLevel <= DebugLog {
		fmt.Printf("debu] "+format+"\n", args...)
	}
}

func pinfo(args ...interface{}) {
	if LogLevel <= InfoLog {
		args = append([]interface{}{"info]"}, args...)
		fmt.Println(args...)
	}
}

func pinfof(format string, args ...interface{}) {
	if LogLevel <= InfoLog {
		fmt.Printf("info] "+format+"\n", args...)
	}
}

func perror(args ...interface{}) {
	if LogLevel <= ErrorLog {
		args = append([]interface{}{"erro]"}, args...)
		fmt.Println(args...)
	}
}

func perrorf(format string, args ...interface{}) {
	if LogLevel <= ErrorLog {
		fmt.Printf("erro] "+format+"\n", args...)
	}
}
