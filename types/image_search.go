package types

type (
	// ImageSearchResult represents custom search image result
	ImageSearchResult struct {
		URL          string
		ThumbnailURL string
	}

	// ImageSearcher represents google custom search image api
	ImageSearcher interface {
		Search(string) (ImageSearchResult, error)
	}
)
