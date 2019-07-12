package problem

import (
	"fmt"
	"io"
	"strconv"
	"strings"
)

// Cause ...
type Cause struct {
	Message    string
	StackTrace StackTrace
}

// Format ...
func (c *Cause) Format(s fmt.State, verb rune) {
	switch verb {
	case 'v':
		if s.Flag('+') {
			io.WriteString(s, c.Message)
			if c.StackTrace != nil {
				c.StackTrace.Format(s, verb)
			}
			return
		}
		fallthrough
	case 's':
		io.WriteString(s, c.Message)
	case 'q':
		fmt.Fprintf(s, "%q", c.Message)
	}
}

// Causes ...
type Causes []Cause

// Format ...
func (c Causes) Format(s fmt.State, verb rune) {
	for _, e := range c {
		e.Format(s, verb)
	}
}

// ToMap ...
func (c Causes) ToMap() []map[string]interface{} {
	var arr = make([]map[string]interface{}, 0)
	for _, e := range c {
		arr = append(arr, map[string]interface{}{
			"message":    e.Message,
			"stacktrace": e.StackTrace.ToMap(),
		})
	}
	return arr
}

func errorToCauses(err error) Causes {
	var (
		frm    *StackTraceFrame
		add    StackTraceFrame
		c      *Cause
		causes = make(Causes, 0)
	)
	for _, s := range strings.Split(fmt.Sprintf("%+v", err), "\n") {
		switch {
		case frm != nil && c != nil && strings.HasPrefix(s, "\t"):
			add, frm = *frm, new(StackTraceFrame)
			add.FileWithLine = strings.TrimPrefix(s, "\t")
			l := strings.Split(add.FileWithLine, ":")
			if len(l) > 0 {
				add.File = l[0]
			}
			if len(l) > 1 {
				add.Line, _ = strconv.Atoi(l[1])
			}
			c.StackTrace = append(c.StackTrace, add)
		case frm != nil && frm.FuncName == "":
			frm.FuncName = s
		default:
			var msg string
			if frm != nil && frm.FuncName != "" {
				msg = frm.FuncName
				frm = &StackTraceFrame{FuncName: s}
			} else {
				msg = s
				frm = new(StackTraceFrame)
			}
			causes = append(causes, Cause{
				Message:    msg,
				StackTrace: make(StackTrace, 0),
			})
			c = &causes[len(causes)-1]
		}
	}
	return causes
}
