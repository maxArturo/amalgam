package worker

import (
	"reflect"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/maxArturo/amalgam"
)

type mockFetcher struct {
	fetcherCount int
	pending      chan *source
	done         chan *source
	updated      chan []amalgam.Linker
}

func (f *mockFetcher) spawnFetchers(count int, pending chan *source, done chan *source, updated chan []amalgam.Linker) {
	f.fetcherCount = count
	f.pending = pending
	f.done = done
	f.updated = updated
}

type mockSleeper struct {
	pending chan *source
	done    chan *source
}

func (f *mockSleeper) sleepSources(done chan *source, pending chan *source, duration time.Duration) {
	f.done = done
	f.pending = pending
}

func TestFetchJob_Start(t *testing.T) {

	type fields struct {
		fetchInterval int
		numFetchers   int
		fetcher       fetcher
		sleeper       sleeper
	}
	type args struct {
		providers []amalgam.Provider
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   chan []amalgam.Linker
	}{
		{name: "passes correct number of fetchers",
			fields: fields{
				numFetchers: 18,
				fetcher:     &mockFetcher{},
				sleeper:     &mockSleeper{},
			}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := &FetchJob{
				fetchInterval: tt.fields.fetchInterval,
				numFetchers:   tt.fields.numFetchers,
				fetcher:       tt.fields.fetcher,
				sleeper:       tt.fields.sleeper,
			}
			// f.Start(tt.args.providers); !reflect.DeepEqual(got, tt.want) {
			// 	t.Errorf("FetchJob.Start() = %v, want %v", got, tt.want)
			// }

			f.Start(tt.args.providers)

			assert.Equal(t, tt.fields.numFetchers, tt.fields.fetcher.(*mockFetcher).fetcherCount)
			assert.True(t, reflect.DeepEqual(tt.fields.fetcher.(*mockFetcher).pending, tt.fields.sleeper.(*mockSleeper).pending))

		})
	}
}
