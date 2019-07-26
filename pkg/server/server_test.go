package server

import (
	"net/http"
	"testing"

	"github.com/maxArturo/amalgam"
	"github.com/stretchr/testify/mock"
)

func TestServer_Run(t *testing.T) {
	updatedChan := make(chan *[]amalgam.Linker)
	handlerFunc := func(w http.ResponseWriter, r *http.Request) {}

	type fields struct {
		fetcher
		layoutHandler
		portResolver
		defaultProviders *[]amalgam.Provider
		httpServer
		logger
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
		{
			name: "default providers",
			fields: fields{
				fetcher: &mockFetcher{
					updated: updatedChan,
				},
				layoutHandler: &mockLayoutHandler{},
				portResolver:  &mockPortResolver{},
				defaultProviders: &[]amalgam.Provider{
					mockProvider{},
				},
				httpServer: &mockHTTPServer{},
				logger:     &mockLogger{},
			},
			args: args{
				port: ":8080",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			s := &Server{
				fetcher:          tt.fields.fetcher,
				layoutHandler:    tt.fields.layoutHandler,
				portResolver:     tt.fields.portResolver,
				defaultProviders: tt.fields.defaultProviders,
				httpServer:       tt.fields.httpServer,
				logger:           tt.fields.logger,
			}

			// expectations
			tt.fields.portResolver.(*mockPortResolver).On("ResolveAddress", tt.args.port)

			if tt.args.sources == nil {
				tt.fields.fetcher.(*mockFetcher).On("Start", tt.fields.defaultProviders).Return(updatedChan)
			} else {
				tt.fields.fetcher.(*mockFetcher).On("Start", &tt.args.sources).Return(updatedChan)
			}

			tt.fields.layoutHandler.(*mockLayoutHandler).On("newHandler", updatedChan).Return(handlerFunc)
			tt.fields.httpServer.(*mockHTTPServer).On("HandleFunc", "/", mock.AnythingOfType("func(http.ResponseWriter, *http.Request)"))
			tt.fields.httpServer.(*mockHTTPServer).On("ListenAndServe", tt.args.port, nil).Return(nil)

			s.Run(tt.args.port, tt.args.sources...)
		})
	}
}
