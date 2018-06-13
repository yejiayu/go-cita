package errors

import "fmt"

type Error interface {
	Error() string
	Build(format string, params ...interface{}) Error
}

func new(errString string) Error {
	return &msg{
		errString: errString,
	}
}

type msg struct {
	errString string
}

func (m *msg) Error() string {
	return m.errString
}

func (m *msg) Build(format string, params ...interface{}) Error {
	message := fmt.Sprintf(format, params)
	m.errString = fmt.Sprintf("%s, %s", m.errString, message)
	return m
}
