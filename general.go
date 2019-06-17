package problem

import (
	"fmt"
	"strings"
)

var (
	// ErrBadRequest is a generic bad request error.
	ErrBadRequest = Type("bad_request", 400)

	// ErrUnauthorized is a generic unauthorized error.
	ErrUnauthorized = Type("unauthorized", 401)

	// ErrInvalidRequest is the Typed of errors produced by the generated code when a request
	// parameter or payload fails to validate.
	ErrInvalidRequest = Type("invalid_request", 400)

	// ErrInvalidEncoding is the error produced when a request body fails to be decoded.
	ErrInvalidEncoding = Type("invalid_encoding", 400)

	// ErrRequestBodyTooLarge is the error produced when the size of a request body exceeds
	// MaxRequestBodyLength bytes.
	ErrRequestBodyTooLarge = Type("request_too_large", 413)

	// ErrNoAuthMiddleware is the error produced when no auth middleware is mounted for a
	// security scheme Typed in the design.
	ErrNoAuthMiddleware = Type("no_auth_middleware", 500)

	// ErrInvalidFile is the error produced by ServeFiles when requested to serve non-existant
	// or non-readable files.
	ErrInvalidFile = Type("invalid_file", 404)

	// ErrNotFound is the error returned to requests that don't match a registered handler.
	ErrNotFound = Type("not_found", 404)

	// ErrMethodNotAllowed is the error returned to requests that match the path of a registered
	// handler but not the HTTP method.
	ErrMethodNotAllowed = Type("method_not_allowed", 405)

	// ErrIO is the Typed of io error
	ErrIO = Type("io_error", 500, Detail("io error"))

	// ErrInternal is the Typed of error used for uncaught errors.
	ErrInternal = Type("internal_error", 500)
)

// OK ...
func OK(s string) error {
	return Type("ok", 200)(s)
}

// MissingPayloadError is the error produced when a request is missing a required payload.
func MissingPayloadError() error {
	return Type("missing_payload", 400, Category("general"))(
		"missing required payload",
	)
}

// InvalidParamTypeError is the error produced when the type of a parameter does not match the type
// Typed in the design.
func InvalidParamTypeError(name string, val interface{}, expected string) error {
	return Type("invalid_param_type", 400, Category("general"))(
		fmt.Sprintf("invalid value %#v for parameter %#v, must be a %s", val, name, expected),
		ValidationResult{
			"name":     name,
			"value":    val,
			"expected": expected,
		},
	)
}

// MissingParamError is the error produced for requests that are missing path or querystring
// parameters.
func MissingParamError(name string) error {
	return Type("missing_param", 400, Category("general"))(
		fmt.Sprintf("missing required parameter %#v", name),
		ValidationResult{
			"name": name,
		},
	)
}

// InvalidAttributeTypeError is the error produced when the type of payload field does not match
// the type Typed in the design.
func InvalidAttributeTypeError(ctx string, val interface{}, expected string) error {
	return Type("invalid_attribute_type", 400, Category("general"))(
		fmt.Sprintf("type of %s must be %s but got value %#v", ctx, expected, val),
		ValidationResult{
			"name":     ctx,
			"value":    val,
			"expected": expected,
		},
	)
}

// MissingAttributeError is the error produced when a request payload is missing a required field.
func MissingAttributeError(ctx, name string) error {
	return Type("missing_attribute", 400, Category("general"))(
		fmt.Sprintf("attribute %#v of %s is missing and required", name, ctx),
		ValidationResult{
			"name":   name,
			"parent": ctx,
		},
	)
}

// MissingHeaderError is the error produced when a request is missing a required header.
func MissingHeaderError(name string) error {
	return Type("missing_header", 400, Category("general"))(
		fmt.Sprintf("missing required HTTP header %#v", name),
		ValidationResult{
			"name": name,
		},
	)
}

// InvalidEnumValueError is the error produced when the value of a parameter or payload field does
// not match one the values Typed in the design Enum validation.
func InvalidEnumValueError(ctx string, val interface{}, allowed []interface{}) error {
	elems := make([]string, len(allowed))
	for i, a := range allowed {
		elems[i] = fmt.Sprintf("%#v", a)
	}
	return Type("invalid_enum_value", 400, Category("general"))(
		fmt.Sprintf("value of %s must be one of %s but got value %#v", ctx, strings.Join(elems, ", "), val),
		ValidationResult{
			"name":  ctx,
			"value": val,
			"expected": strings.Join(elems, ": "),
		},
	)
}

// InvalidFormatError is the error produced when the value of a parameter or payload field does not
// match the format validation Typed in the design.
func InvalidFormatError(ctx, target string, format string, formatError error) error {
	return Type("invalid_format", 400, Category("general"))(
		fmt.Sprintf("%s must be formatted as a %s but got value %#v, %s", ctx, format, target, formatError.Error()),
		ValidationResult{
			"name":     ctx,
			"value":    target,
			"expected": format,
			"error":    formatError.Error(),
		},
	)
}

// InvalidPatternError is the error produced when the value of a parameter or payload field does
// not match the pattern validation Typed in the design.
func InvalidPatternError(ctx, target string, pattern string) error {
	return Type("invalid_pattern", 400, Category("general"))(
		fmt.Sprintf("%s must match the regexp %#v but got value %#v", ctx, pattern, target),
		ValidationResult{
			"name":     ctx,
			"value":    target,
			"expected": pattern,
		},
	)
}

// InvalidRangeError is the error produced when the value of a parameter or payload field does
// not match the range validation Typed in the design. value may be a int or a float64.
func InvalidRangeError(ctx string, target interface{}, value interface{}, min bool) error {
	comp := "greater than or equal to"
	if !min {
		comp = "less than or equal to"
	}
	return Type("invalid_range", 400, Category("general"))(
		fmt.Sprintf("%s must be %s %v but got value %#v", ctx, comp, value, target),
		ValidationResult{
			"name":     ctx,
			"value":    target,
			"comp":     comp,
			"expected": value,
			"min":      min,
		},
	)
}

// InvalidLengthError is the error produced when the value of a parameter or payload field does
// not match the length validation Typed in the design.
func InvalidLengthError(ctx string, target interface{}, ln, value int, min bool) error {
	comp := "greater than or equal to"
	if !min {
		comp = "less than or equal to"
	}
	return Type("invalid_length", 400, Category("general"))(
		fmt.Sprintf("length of %s must be %s %d but got value %#v (len=%d)", ctx, comp, value, target, ln),
		ValidationResult{
			"name":     ctx,
			"value":    target,
			"len":      ln,
			"comp":     comp,
			"expected": value,
			"min":      min,
		},
	)
}

// NoAuthMiddleware is the error produced when goa is unable to lookup a auth middleware for a
// security scheme Typed in the design.
func NoAuthMiddleware(schemeName string) error {
	return Type("no_auth_middleware", 500, Category("general"))(
		fmt.Sprintf("Auth middleware for security scheme %s is not mounted", schemeName),
		ValidationResult{
			"scheme": schemeName,
		},
	)
}

// MethodNotAllowedError is the error produced to requests that match the path of a registered
// handler but not the HTTP method.
func MethodNotAllowedError(method string, allowed []string) error {
	var plural string
	if len(allowed) > 1 {
		plural = " one of"
	}
	return Type("method_not_allowed", 405, Category("general"))(
		fmt.Sprintf("Method %s must be%s %s", method, plural, strings.Join(allowed, ", ")),
		ValidationResult{
			"method": method,
			"allowed": strings.Join(allowed, ": "),
		},
	)
}
