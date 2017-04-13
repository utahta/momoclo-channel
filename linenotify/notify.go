package linenotify

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/url"
	"path"
	"strings"
	"sync"

	"github.com/pkg/errors"
	"github.com/utahta/go-openuri"
)

type RequestNotify struct {
	Client         *http.Client
	imageBodyCache map[string][]byte
	mux            *sync.Mutex
}

var (
	ErrorNotifyInvalidAccessToken = errors.New("Invalid access token.")
)

func NewRequestNotify() *RequestNotify {
	return &RequestNotify{Client: http.DefaultClient, imageBodyCache: map[string][]byte{}, mux: new(sync.Mutex)}
}

func (r *RequestNotify) Notify(token, message, imageThumbnail, imageFullsize, imageFile string) error {
	var (
		contentType string
		body        io.Reader
		err         error
	)

	if imageFile != "" {
		if body, contentType, err = r.requestBodyWithImageFile(message, imageFile); err != nil {
			return err
		}
	} else {
		if body, contentType, err = r.requestBody(message, imageThumbnail, imageFullsize); err != nil {
			return err
		}
	}

	req, err := http.NewRequest("POST", "https://notify-api.line.me/api/notify", body)
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", contentType)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))

	resp, err := r.Client.Do(req)
	if err != nil {
		return err
	}

	if resp.StatusCode == http.StatusUnauthorized {
		return ErrorNotifyInvalidAccessToken
	}

	if resp.StatusCode != http.StatusOK {
		content, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return err
		}

		var data interface{}
		err = json.Unmarshal(content, &data)
		if err != nil {
			return err
		}
		root := data.(map[string]interface{})
		return errors.New(root["message"].(string))
	}
	return nil
}

func (r *RequestNotify) requestBody(message, imageThumbnail, imageFullsize string) (io.Reader, string, error) {
	v := url.Values{}
	v.Add("message", message)
	if imageThumbnail != "" {
		v.Add("imageThumbnail", imageThumbnail)
	}
	if imageFullsize != "" {
		v.Add("imageFullsize", imageFullsize)
	}
	return strings.NewReader(v.Encode()), "application/x-www-form-urlencoded", nil
}

func (r *RequestNotify) requestBodyWithImageFile(message, imageFile string) (io.Reader, string, error) {
	var b bytes.Buffer
	w := multipart.NewWriter(&b)

	if err := w.WriteField("message", message); err != nil {
		return nil, "", err
	}

	fw, err := w.CreateFormFile("imageFile", path.Base(imageFile))
	if err != nil {
		return nil, "", err
	}

	cache, err := r.CacheImage(imageFile)
	if err != nil {
		return nil, "", err
	}

	if _, err := io.Copy(fw, bytes.NewReader(cache)); err != nil {
		return nil, "", err
	}

	if err := w.Close(); err != nil {
		return nil, "", err
	}

	return &b, w.FormDataContentType(), nil
}

// Cache image, returns bytes of image
func (r *RequestNotify) CacheImage(imageFile string) ([]byte, error) {
	r.mux.Lock()
	defer r.mux.Unlock()

	if cache, ok := r.imageBodyCache[imageFile]; ok {
		return cache, nil
	}

	o, err := openuri.Open(imageFile, openuri.WithHTTPClient(r.Client))
	if err != nil {
		return nil, err
	}
	defer o.Close()

	imgSrc, err := ioutil.ReadAll(o)
	if err != nil {
		return nil, err
	}
	r.imageBodyCache[imageFile] = make([]byte, len(imgSrc))
	copy(r.imageBodyCache[imageFile], imgSrc)

	return r.imageBodyCache[imageFile], nil
}

func (r *RequestNotify) ClearImage(imageFile string) {
	delete(r.imageBodyCache, imageFile)
}
