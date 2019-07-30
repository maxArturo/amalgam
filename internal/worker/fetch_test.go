package worker

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/maxArturo/amalgam"
)

func Test_fetchProvider_fetch(t *testing.T) {
	tests := []struct {
		name        string
		f           *SourceFetch
		source      *source
		sourceFails bool
	}{
		{
			name: "fetches source normally",
			f:    &SourceFetch{},
			source: &source{
				provider: &mockProvider{},
			},
		},
		{
			name: "adds error count when fetch fails",
			f:    &SourceFetch{},
			source: &source{
				provider: &mockProvider{
					returnFetchErr: true,
				},
			},
			sourceFails: true,
		},
	}
	for i, tt := range tests {
		in := make(chan *source)
		out := make(chan *source)
		content := make(chan *[]amalgam.Linker)

		t.Run(tt.name, func(t *testing.T) {
			go tt.f.fetch(i, in, out, content)
			in <- tt.source

			if tt.sourceFails {
				outSrc := <-out
				assert.True(t, tt.source.errCount > 0)
				assert.Equal(t, tt.source, outSrc)
			} else {
				links := <-content
				assert.NotNil(t, links)
				assert.Equal(t, 0, tt.source.errCount)
			}
		})
	}
}
