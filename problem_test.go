package problem

import (
	"fmt"
	"testing"
)

func TestNew_get_error_message(t *testing.T) {
	typ := ErrInternal("test")
	if e, a := "internal_error: test", typ.Error(); e != a {
		t.Errorf("{\n - %v\n + %v\n}", e, a)
	}
}

func TestNew_format_s(t *testing.T) {
	err := ErrInternal("test")
	if e, a := "internal_error: test", fmt.Sprintf("%s", err); e != a {
		t.Errorf("{\n - %v\n + %v\n}", e, a)
	}
}

func TestNew_format_v(t *testing.T) {
	err := ErrInternal("test")
	if e, a := "internal_error: test", fmt.Sprintf("%v", err); e != a {
		t.Errorf("{\n - %v\n + %v\n}", e, a)
	}
}

// func TestNew_format_plus_v(t *testing.T) {
// 	err := ErrInternal("test")
// 	if e, a := "internal_server_error: test", fmt.Sprintf("%+v", err); e != a {
// 		t.Errorf("{\n - %v\n + %v\n}", e, a)
// 	}
// }

// func TestNew_with_stack(t *testing.T) {
// 	src := ErrInternal("test")
// 	err := ErrInternal(src)
// 	if e, a := "internal_server_error: test", fmt.Sprintf("%+v", err); e != a {
// 		t.Errorf("{\n - %v\n + %v\n}", e, a)
// 	}
// }
