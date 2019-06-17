package problem

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"io"
)

const (
	// ErrorTypeBaseURI ...
	ErrorTypeBaseURI      = ""
	defaultErrTypeBaseURI = "https://example.com/errors"
	defaultErrType        = "about:blank"
)

// Typed ...
type Typed func(vals ...interface{}) error

// MergeableError ...
type MergeableError interface {
	error
	Merge(err error) error
}

// Detail ...
type Detail string

// Category ...
type Category string

// Instance ...
type Instance string

// Resource ...
type Resource string

// Field ...
type Field string

// Meta ...
type Meta []interface{}

func (l Meta) mergeTo(m map[string]interface{}) map[string]interface{} {
	if m == nil {
		m = make(map[string]interface{})
	}
	var k string
	for n, v := range l {
		if n%2 == 0 {
			k = fmt.Sprint(v)
		} else {
			m[k] = v
		}
	}
	return m
}

// ValidationResult ...
type ValidationResult map[string]interface{}

func (l ValidationResult) mergeTo(m map[string]interface{}) map[string]interface{} {
	if m == nil {
		m = make(map[string]interface{})
	}
	for k, v := range l {
		m[k] = v
	}
	return m
}

// Type ...
func Type(title string, status int, typeVals ...interface{}) Typed {
	return func(vals ...interface{}) error {
		var (
			c = Cause{StackTrace: Callers()}
			p = &Problem{
				ID:     problemID(),
				Title:  title,
				Status: status,
				Causes: Causes{c},
			}
		)
		// append values
		if len(typeVals) > 0 {
			vals = append(typeVals, vals...)
		}
		for i, any := range vals {
			switch any := any.(type) {
			case Detail:
				p.Detail = string(any)
			case Category:
				p.Category = string(any)
			case Instance:
				p.Instance = string(any)
			case Resource:
				p.Resource = string(any)
			case Field:
				p.Field = string(any)
			case string:
				if p.Detail == "" {
					p.Detail = any
				}
			case error:
				c.Message = any.Error()
				p.err = any
				p.Causes = append(errorToCauses(any), p.Causes...)
				if p.Detail == "" {
					p.Detail = any.Error()
				}
			case fmt.Stringer:
				if p.Detail == "" {
					p.Detail = any.String()
				}
			case Meta:
				p.Meta = any.mergeTo(p.Meta)
			case ValidationResult:
				p.ValidationResult = any.mergeTo(p.ValidationResult)
			default:
				if i == 0 {
					p.Detail = fmt.Sprintf("%v", any)
				}
			}
		}
		if p.Detail == "" {
			p.Detail = "unknown"
		}
		if c.Message == "" {
			c.Message = p.Detail
		}
		p.Type = typeURI(p.Category, p.Title)
		return p
	}
}

// WithinCategory ...
func WithinCategory(category, title string, status int, typeVals ...interface{}) Typed {
	var add = make([]interface{}, 0, len(typeVals)+1)
	add = append(add, Category(category))
	if len(typeVals) > 0 {
		add = append(add, typeVals...)
	}
	return Type(title, status, add...)
}

// Problem ...
type Problem struct {
	err              error
	ID               string
	Type             string
	Category         string
	Title            string
	Status           int
	Detail           string
	Instance         string
	Resource         string
	Field            string
	Merged           []error
	Meta             map[string]interface{}
	Causes           Causes
	ValidationResult map[string]interface{}
}

// Errror implements error.Error
func (p *Problem) Error() string {
	if p.Detail != "" {
		return fmt.Sprintf("%s: %s", p.Title, p.Detail)
	}
	return fmt.Sprintf("%s", p.Title)
}

// ResponseStatus implements github.com/goadesign/goa/ServiceError.ResponseStatus
func (p *Problem) ResponseStatus() int {
	return p.Status
}

// Token implements github.com/goadesign/goa/ServiceError.Token
func (p *Problem) Token() string {
	return p.ID
}

// Merge implements github.com/goadesign/goa/ServiceMergeableError.Merge
func (p *Problem) Merge(err error) error {
	if err == nil {
		return p
	}

	var (
		src = *p
		o   = *AsProblem(err)
		m   = Problem{
			ID:       problemID(),
			Type:     p.Type,
			Category: p.Category,
			Title:    p.Title,
			Status:   p.Status,
			Detail:   p.Detail,
			Meta:     map[string]interface{}{"merged": true},
		}
	)

	if o.Detail != "" {
		m.Detail = fmt.Sprintf("%s; %s", m.Detail, o.Detail)
	}

	var (
		errsX = src.Merged
		errsY = o.Merged
	)
	src.Merged = nil
	o.Merged = nil
	if len(errsX) == 0 {
		errsX = []error{&src}
	}
	if len(errsY) == 0 {
		errsY = []error{&o}
	}
	m.Merged = append(errsX, errsY...)

	if src.Status != o.Status || src.Title != o.Title {
		if src.Status < o.Status {
			m.Status = o.Status
		}
		switch {
		case m.Status < 500:
			m.Status = 400
			m.Title = "bad_request"
		default:
			m.Status = 500
			m.Title = "internal_error"
		}
		m.Category = "merged"
		m.Type = typeURI(m.Category, m.Title)
	}

	return &m
}

// Format ...
func (p *Problem) Format(s fmt.State, verb rune) {
	msg := p.Error()
	switch verb {
	case 'v':
		if s.Flag('+') {
			p.Causes.Format(s, verb)
			return
		}
		fallthrough
	case 's':
		io.WriteString(s, msg)
	case 'q':
		fmt.Fprintf(s, "%q", msg)
	}
}

// AsProblem ...
func AsProblem(err error) *Problem {
	p, ok := err.(*Problem)
	if !ok {
		p, _ = ErrInternal(err).(*Problem)
	}
	return p
}

func problemID() string {
	b := make([]byte, 6)
	io.ReadFull(rand.Reader, b)
	return base64.StdEncoding.EncodeToString(b)
}

func typeURI(category, title string) string {
	if category != "" && title != "" {
		baseErrType := defaultErrTypeBaseURI
		if s := ErrorTypeBaseURI; s != "" {
			baseErrType = s
		}
		return fmt.Sprintf("%s/%s/%s", baseErrType, category, title)
	}
	return defaultErrType
}
