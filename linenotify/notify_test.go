package linenotify

import (
	"io/ioutil"
	"net/http"
	"strings"
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
		{&http.Response{StatusCode: http.StatusOK}, nil},
		{&http.Response{StatusCode: http.StatusUnauthorized}, ErrorNotifyInvalidAccessToken},
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
		{2, &http.Response{}, errors.New("Not cached!!")},
	}

	for _, test := range tests {
		req.Client.Transport = &notifyRoundTripper{
			resp: test.resp,
			err:  test.err,
		}

		body, contentType, err := req.requestBodyWithImageFile("test", "dummy.jpg")
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
}
