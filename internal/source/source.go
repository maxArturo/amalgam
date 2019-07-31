package source

import (
	"time"

	"github.com/maxArturo/amalgam"
)

// Source is the canonical internal representation of a Provider.
type Source struct {
	Provider    amalgam.Provider
	ErrCount    int
	LastUpdated time.Time
}
