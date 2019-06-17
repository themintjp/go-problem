package problem

import (
	"fmt"
	"runtime"
)

// StackTrace ...
type StackTrace []StackTraceFrame

// Shift ...
func (st StackTrace) Shift() StackTrace {
	if len(st) > 1 {
		return st[1:]
	}
	return st
}

// Format ...
func (st StackTrace) Format(s fmt.State, verb rune) {
	for _, f := range st {
		f.Format(s, verb)
	}
}

// ToMap ...
func (st StackTrace) ToMap() []map[string]string {
	var arr = make([]map[string]string, 0)
	for _, e := range st {
		arr = append(arr, map[string]string{
			"file": e.fileWithLine,
			"func": e.funcName,
		})
	}
	return arr
}

// StackTraceFrame ...
type StackTraceFrame struct {
	file,
	fileWithLine,
	funcName string
	line int
}

// Format ...
func (f *StackTraceFrame) Format(s fmt.State, verb rune) {
	if verb == 'v' && s.Flag('+') {
		fmt.Fprintf(s, "\n%s\n\t%s", f.funcName, f.fileWithLine)
	}
}

// Callers ...
func Callers() StackTrace {
	const depth = 32
	var (
		pcs [depth]uintptr
		n = runtime.Callers(3, pcs[:])
		st = make(StackTrace, n)
	)
	for i := 0; i < n; i++ {
		var (
			pc = pcs[i] -1
			f = &st[i]
		)
		f.funcName = "unknown"
		f.file = "unknown"
		f.line = 0
		if fn := runtime.FuncForPC(pc); fn != nil {
			f.file, f.line = fn.FileLine(pc)
			f.fileWithLine = fmt.Sprintf("%s:%d", f.file, f.line)
			f.funcName = fn.Name()
		}
	}
	return st
}