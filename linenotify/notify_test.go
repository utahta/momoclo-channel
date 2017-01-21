package linenotify

import (
	"testing"
	"net/http"
)

type notifyRoundTripper struct {
	resp *http.Response
}

func (rt *notifyRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	rt.resp.Request = req
	return rt.resp, nil
}

func TestRequestNotify_Notify(t *testing.T) {
	req := NewRequestNotify()
	tests := []struct{
		resp *http.Response
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
