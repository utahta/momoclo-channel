# momoclo-crawler

Momoiro Clover Z's blogs, ae news, google news and youtube videos crawler.

# Install

```
$ go get github.com/utahta/momoclo-channel/crawler
```

# Usage

```go
package main

import (
    "log"
    "github.com/utahta/momoclo-channel/crawler"
)

func main() {
    items, err := crawler.FetchAeNews()
    if err != nil {
        log.Fatal(err)
    }

    for _, item := range items {
        log.Println(item.Url, item.Title, item.PublishedAt)
	}
}
```
