package utils

import (
	"runtime"
	"strings"
)

func GetFunctionName(skip int) string {
	pc, _, _, ok := runtime.Caller(skip)
	if !ok {
		return "unknown"
	}
	fn := runtime.FuncForPC(pc)
	if fn == nil {
		return "unknown"
	}
	// Remove full path and keep just package.FunctionName
	parts := strings.Split(fn.Name(), "/")
	return parts[len(parts)-1]
}
