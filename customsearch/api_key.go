package customsearch

type (
	apiKey string
)

func (a apiKey) Get() (string, string) {
	return "key", string(a)
}
