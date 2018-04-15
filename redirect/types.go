package redirect

// Redirect represents a mapping from a key to a URL.
type Redirect struct {
	Key string
	URL string
}

// Redirector allows fetching a Redirect for given key.
type Redirector interface {
	Get(key string) (*Redirect, error)
}
