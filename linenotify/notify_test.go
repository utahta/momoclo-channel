package linenotify

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"sync"
	"testing"

	"github.com/pkg/errors"
)

type notifyRoundTripper struct {
	resp *http.Response
	err  error
}

func (rt *notifyRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	rt.resp.Request = req
	return rt.resp, rt.err
}

func TestRequestNotify_Notify(t *testing.T) {
	req := NewRequestNotify()
	tests := []struct {
		resp        *http.Response
		expectedErr error
	}{
		{&http.Response{StatusCode: http.StatusOK, Body: ioutil.NopCloser(strings.NewReader(""))}, nil},
		{&http.Response{StatusCode: http.StatusUnauthorized, Body: ioutil.NopCloser(strings.NewReader(""))}, ErrorNotifyInvalidAccessToken},
	}

	for _, test := range tests {
		req.Client.Transport = &notifyRoundTripper{resp: test.resp}

		err := req.Notify("token", "test", "", "", "")
		if err != test.expectedErr {
			t.Fatal(err)
		}
	}
}

func TestNewRequestNotify_requestBodyWithImageFile(t *testing.T) {
	req := NewRequestNotify()

	tests := []struct {
		idx  int
		resp *http.Response
		err  error
	}{
		{1, &http.Response{StatusCode: http.StatusOK, Body: ioutil.NopCloser(strings.NewReader("image file"))}, nil},
		{2, &http.Response{Body: ioutil.NopCloser(strings.NewReader(""))}, errors.New("Not cached!!")},
	}

	for _, test := range tests {
		req.Client.Transport = &notifyRoundTripper{resp: test.resp, err: test.err}

		body, contentType, err := req.requestBodyWithImageFile("test", "http://localhost/dummy.jpg")
		if err != nil {
			t.Fatal(err)
		}
		buf, err := ioutil.ReadAll(body)
		if err != nil {
			t.Fatal(err)
		}

		if !strings.Contains(string(buf), "image file") {
			t.Errorf("Expected buffer image file[%d], got %s", test.idx, string(buf))
		}

		if !strings.Contains(contentType, "multipart/form-data;") {
			t.Errorf("Expected contentType[%d], got %s", test.idx, contentType)
		}
	}

	// ファイル名を変えた場合、cache が無効になり GET すること
	req.Client.Transport = &notifyRoundTripper{
		resp: &http.Response{},
		err:  errors.New("expect call"),
	}
	_, _, err := req.requestBodyWithImageFile("test", "http://localhost/dummy2.jpg")
	if err.Error() != "Get http://localhost/dummy2.jpg: expect call" {
		t.Fatalf("Expected error, got %v", err)
	}

	// test for data race
	var wg sync.WaitGroup
	req.Client.Transport = &notifyRoundTripper{resp: &http.Response{Body: ioutil.NopCloser(strings.NewReader("image file"))}, err: nil}
	for _, test := range tests {
		wg.Add(1)
		go func(idx int) {
			defer wg.Done()
			_, _, err := req.requestBodyWithImageFile("test", fmt.Sprintf("http://localhost/dummy_%d.jpg)", idx))
			if err != nil {
				t.Fatal(err)
			}
		}(test.idx)
	}
	wg.Wait()

	for i := 0; i < 4; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			_, _, err := req.requestBodyWithImageFile("test", "http://localhost/dummy.jpg)")
			if err != nil {
				t.Fatal(err)
			}
		}()
	}
	wg.Wait()
}
