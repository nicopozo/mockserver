package mocks

import (
	"bufio"
	"net"
	"net/http"
)

type MockGinResponseWriter struct {
	http.ResponseWriter
	http.Hijacker
	http.Flusher
	http.CloseNotifier
	Bytes  []byte
	status int
}

func (mr *MockGinResponseWriter) Status() int {
	return mr.status
}

func (mr *MockGinResponseWriter) Size() int {
	return int(0)
}

func (mr *MockGinResponseWriter) WriteString(s string) (int, error) {
	mr.Bytes = []byte(s)

	return 0, nil
}

func (mr *MockGinResponseWriter) Written() bool {
	return false
}

func (mr *MockGinResponseWriter) WriteHeaderNow() {
}

func (mr *MockGinResponseWriter) Header() http.Header {
	return make(map[string][]string)
}

func (mr *MockGinResponseWriter) Write(b []byte) (int, error) {
	mr.Bytes = b

	return 0, nil
}

func (mr *MockGinResponseWriter) WriteHeader(status int) {
	mr.status = status
}

func (mr *MockGinResponseWriter) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	return nil, nil, nil
}

func (mr *MockGinResponseWriter) Flush() {
}

func (mr *MockGinResponseWriter) CloseNotify() <-chan bool {
	return nil
}

//nolint:nonamedreturns
func (mr *MockGinResponseWriter) Pusher() (pusher http.Pusher) {
	return nil
}
