package fetch

// Fetcher is the interface that wraps the Fetch method.
// A Fetcher should return the same resource with every fetch, and `Fetch` should only return different data when the resource changes.
//
// Fetch returns bytes fetched from some source, and any error.
type Fetcher interface {
	Fetch() ([]byte, error)
}
