package model

type (
	// CustomSearchImageResult represents custom search image result
	CustomSearchImageResult struct {
		URL          string
		ThumbnailURL string
	}

	// CustomSearchClient represents google custom search api
	CustomSearchClient interface {
		SearchImage(string) (CustomSearchImageResult, error)
	}
)
