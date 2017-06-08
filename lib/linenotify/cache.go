package linenotify

import (
	"bytes"
	"io"
	"net/http"
	"sync"

	"github.com/utahta/go-openuri"
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

	o, err := openuri.Open(filename, openuri.WithHTTPClient(c))
	if err != nil {
		return nil, err
	}
	defer o.Close()

	buf := &bytes.Buffer{}
	if _, err := io.Copy(buf, o); err != nil {
		return nil, err
	}
	images[filename] = buf.Bytes()

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
