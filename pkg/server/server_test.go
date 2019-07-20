package server

import (
	"net/http"
	"testing"

	"github.com/maxArturo/amalgam"
	"github.com/stretchr/testify/mock"
)

type mockFetcher struct {
	mock.Mock
}

func (f *mockFetcher) Start(providers *[]amalgam.Provider) chan *[]amalgam.Linker {
	f.Called(providers)
	return make(chan *[]amalgam.Linker)
}

type mockLayoutHandler struct {
	mock.Mock
}

func (f *mockLayoutHandler) newHandler(in chan *[]amalgam.Linker) func(w http.ResponseWriter, r *http.Request) {
	f.Called(in)
	return func(w http.ResponseWriter, r *http.Request) {}
}

type mockPortResolver struct {
	mock.Mock
}

func TestServer_Run(t *testing.T) {
	type fields struct {
		fetcher          fetcher
		layoutHandler    layoutHandler
		portResolver     portResolver
		defaultProviders *[]amalgam.Provider
	}
	type args struct {
		port    string
		sources []amalgam.Provider
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Server{
				fetcher:          tt.fields.fetcher,
				layoutHandler:    tt.fields.layoutHandler,
				portResolver:     tt.fields.portResolver,
				defaultProviders: tt.fields.defaultProviders,
			}
			s.Run(tt.args.port, tt.args.sources...)
		})
	}
}
