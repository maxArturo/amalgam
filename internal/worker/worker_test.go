package worker

import (
	"errors"
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

type mockProvider struct {
	returnFetchErr bool
}

func (f mockProvider) Fetch() ([]amalgam.Linker, error) {
	if f.returnFetchErr {
		return nil, errors.New("mock error")
	}
	return []amalgam.Linker{}, nil
}
func (f mockProvider) Name() string {
	return "Mock Provider"
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
		name         string
		fields       fields
		args         args
		want         chan []amalgam.Linker
		flexChannels bool
	}{
		{
			name: "passes correct number of fetchers",
			fields: fields{
				numFetchers: 18,
				fetcher:     &mockFetcher{},
				sleeper:     &mockSleeper{},
			}}, {
			name: "passes sources to pending channel",
			fields: fields{
				fetcher: &mockFetcher{},
				sleeper: &mockSleeper{},
			},
			args: args{providers: []amalgam.Provider{
				mockProvider{},
				mockProvider{},
				mockProvider{},
			}},
			flexChannels: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := &FetchJob{
				fetchInterval: tt.fields.fetchInterval,
				numFetchers:   tt.fields.numFetchers,
				fetcher:       tt.fields.fetcher,
				sleeper:       tt.fields.sleeper,
			}

			f.Start(tt.args.providers)

			pendingChan := tt.fields.fetcher.(*mockFetcher).pending

			assert.Equal(t, tt.fields.numFetchers, tt.fields.fetcher.(*mockFetcher).fetcherCount)
			assert.True(t, reflect.DeepEqual(tt.fields.fetcher.(*mockFetcher).pending, tt.fields.sleeper.(*mockSleeper).pending))
			assert.True(t, reflect.DeepEqual(tt.fields.fetcher.(*mockFetcher).done, tt.fields.sleeper.(*mockSleeper).done))

			if tt.flexChannels {
				for _, p := range tt.args.providers {
					nextSource := <-pendingChan
					assert.Equal(t, nextSource.provider, p)
				}

			}
		})
	}
}
