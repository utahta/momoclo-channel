package linenotify

import (
	"bytes"
	"io"
	"net/http"
	"strings"
	"sync"

	"github.com/pkg/errors"
	"github.com/utahta/go-openuri"
	"github.com/utahta/momoclo-channel/lib/backoff"
)

var (
	images = map[string][]byte{}
	mux    = &sync.Mutex{}
)

func fetchImage(c *http.Client, filename string) ([]byte, error) {
	mux.Lock()
	defer mux.Unlock()

	if image, ok := images[filename]; ok {
		return image, nil
	}

	err := backoff.Retry(3, func() error {
		o, err := openuri.Open(filename, openuri.WithHTTPClient(c))
		if err != nil {
			return err
		}
		defer o.Close()

		buf := &bytes.Buffer{}
		if _, err := io.Copy(buf, o); err != nil {
			return err
		}

		if ct := http.DetectContentType(buf.Bytes()); !strings.Contains(ct, "image") {
			return errors.Errorf("Detected invalid content type. ct:%v", ct)
		}
		images[filename] = buf.Bytes()

		return nil
	})
	if err != nil {
		return nil, err
	}

	return images[filename], nil
}

func clearImage(filename string) {
	mux.Lock()
	defer mux.Unlock()

	delete(images, filename)
}

func cacheImage(filename string) []byte {
	mux.Lock()
	defer mux.Unlock()

	if image, ok := images[filename]; ok {
		return image
	}
	return nil
}
