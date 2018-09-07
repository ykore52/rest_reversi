package server

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
)

type FakeHTTPResponseWriter struct {
	http.ResponseWriter
	FakeWriteHeader func(statusCode int)
	FakeWrite       func(stream []byte) (int, error)
}

func (f *FakeHTTPResponseWriter) WriteHeader(statusCode int) {
	f.FakeWriteHeader(statusCode)
}

func (f *FakeHTTPResponseWriter) Write(stream []byte) (int, error) {
	return f.FakeWrite(stream)
}

func Test_APIPostUser_Success(t *testing.T) {

	url, _ := url.Parse("http://localhost/user")
	postData := `{"Name": "test1"}`
	r := &http.Request{
		Method:        "POST",
		Body:          ioutil.NopCloser(bytes.NewReader([]byte(postData))),
		ContentLength: int64(len(postData)),
		URL:           url,
		RequestURI:    "/user",
	}

	w := &FakeHTTPResponseWriter{
		FakeWriteHeader: func(statusCode int) { /* nothing to do */ },
		FakeWrite: func(stream []byte) (int, error) {
			assert.Equal(t, string(stream[:19]), `{"status":"success"`)
			return 0, nil
		},
	}

	APIRoute(w, r)
}

func Test_APIPostUser_Error_Method(t *testing.T) {

	url, _ := url.Parse("http://localhost/user")
	postData := `{"Name": "test1"}`
	r := &http.Request{
		Method:        "GET",
		Body:          ioutil.NopCloser(bytes.NewReader([]byte(postData))),
		ContentLength: int64(len(postData)),
		URL:           url,
		RequestURI:    "/user",
	}

	w := &FakeHTTPResponseWriter{
		FakeWriteHeader: func(statusCode int) { /* nothing to do */ },
		FakeWrite: func(stream []byte) (int, error) {
			assert.Equal(t, string(stream[:16]), `{"status":"fail"`)
			return 0, nil
		},
	}

	APIRoute(w, r)
}

func Test_APIPostUser_Error_ContentLength(t *testing.T) {

	url, _ := url.Parse("http://localhost/user")
	postData := `{"Name": "test1"}`
	r := &http.Request{
		Method:        "POST",
		Body:          ioutil.NopCloser(bytes.NewReader([]byte(postData))),
		ContentLength: 0,
		URL:           url,
		RequestURI:    "/user",
	}

	w := &FakeHTTPResponseWriter{
		FakeWriteHeader: func(statusCode int) { /* nothing to do */ },
		FakeWrite: func(stream []byte) (int, error) {
			assert.Equal(t, string(stream[:16]), `{"status":"fail"`)
			return 0, nil
		},
	}

	APIRoute(w, r)
}

func Test_APIPostUser_Error_BrokenJSON(t *testing.T) {

	url, _ := url.Parse("http://localhost/user")
	postData := `"Name": "test1"`
	r := &http.Request{
		Method:        "POST",
		URL:           url,
		Body:          ioutil.NopCloser(bytes.NewReader([]byte(postData))),
		ContentLength: 0,
		RequestURI:    "/user",
	}

	w := &FakeHTTPResponseWriter{
		FakeWriteHeader: func(statusCode int) { /* nothing to do */ },
		FakeWrite: func(stream []byte) (int, error) {
			assert.Equal(t, string(stream[:16]), `{"status":"fail"`)
			return 0, nil
		},
	}

	APIRoute(w, r)
}

func Test_APIGetUser_Success(t *testing.T) {

	done := make(chan GetSessionInfoResponse, 1)

	if true {

		url, _ := url.Parse("http://localhost/user")
		postData := `{"Name": "test1"}`
		r := &http.Request{
			Method:        "POST",
			Body:          ioutil.NopCloser(bytes.NewReader([]byte(postData))),
			ContentLength: int64(len(postData)),
			URL:           url,
			RequestURI:    "/user",
		}

		w := &FakeHTTPResponseWriter{
			FakeWriteHeader: func(statusCode int) { /* nothing to do */ },
			FakeWrite: func(stream []byte) (int, error) {
				assert.Equal(t, string(stream[:19]), `{"status":"success"`)
				s := new(GetSessionInfoResponse)
				json.Unmarshal(stream, s)
				done <- *s
				close(done)
				return 0, nil
			},
		}

		APIRoute(w, r)
	}

	sessionInfo, ok := <-done
	if !ok {
		assert.Fail(t, "Cannot receive data")
		return
	}

	if true {

		url, _ := url.Parse("http://localhost/user/" + sessionInfo.UserID)
		r := &http.Request{
			Method:        "GET",
			URL:           url,
			RequestURI:    "/user",
			ContentLength: 0,
		}

		w := &FakeHTTPResponseWriter{
			FakeWriteHeader: func(statusCode int) { /* nothing to do */ },
			FakeWrite: func(stream []byte) (int, error) {
				fmt.Println("OUTPUT: " + string(stream))
				assert.Equal(t, string(stream[:19]), `{"status":"success"`)
				return 0, nil
			},
		}

		APIRoute(w, r)
	}
}
