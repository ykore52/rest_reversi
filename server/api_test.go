package server

import (
	"net/http"
	"testing"
)

type FakeHttpResponseWriter struct {
	http.ResponseWriter
	FakeWriteHeader func(statusCode int)
	FakeWrite       func(stream []byte) (int, error)
}

func (f *FakeHttpResponseWriter) WriteHeader(statusCode int) {
	f.FakeWriteHeader(statusCode)
}

func (f *FakeHttpResponseWriter) Write(stream []byte) (int, error) {
	return f.FakeWrite(stream)
}

type FakeHttpRequest struct {
	http.Request
	Body          string
	ContentLength int
}

func Test_PostRegist(t *testing.T) {

	w := &FakeHttpResponseWriter{
		FakeWriteHeader: func(statusCode int) {},
		FakeWrite: func(stream []byte) (int, error) {
			return 0, nil
		},
	}

	r := &FakeHttpRequest{
		Body:          "hoge",
		ContentLength: 4,
	}

	PostRegist(w, r)
}
