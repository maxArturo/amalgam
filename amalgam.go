// Package amalgam provides a brutally simple webserver for no-nonsense news aggregation.
// It aims for lightweight hardware requirements, extensibility, and simplicity.
package amalgam

// Provider defines the only requirement from a source: a Fetch() function that returns
// a slice of NewsLink structs, and a Name().
type Provider interface {
	Fetch() ([]Linker, error)
	Name() string
}

// Linker represents a link with a set of values needed to be rendered by the news aggregator.
// These are implemented as methods for flexibility of future providers.
type Linker interface {
	Source() string      // the name of the source
	Title() string       // link title
	URL() string         // link URL
	CommentsURL() string // URL to the comments page (if any)
	CommentCount() int   // count of comments, if available
}
