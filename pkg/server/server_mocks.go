package server

import (
	"net/http"

	"github.com/maxArturo/amalgam"
	"github.com/stretchr/testify/mock"
)

type mockFetcher struct {
	mock.Mock
	updated chan *[]amalgam.Linker
}

func (f *mockFetcher) Start(providers *[]amalgam.Provider) chan *[]amalgam.Linker {
	f.Called(providers)
	return f.updated
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

func (f *mockPortResolver) ResolveAddress(addr string) string {
	return ":8080"
}

type mockHTTPServer struct {
	mock.Mock
}

func (f *mockHTTPServer) HandleFunc(pattern string, handler func(http.ResponseWriter, *http.Request)) {
	f.Called(pattern, handler)
}

func (f *mockHTTPServer) ListenAndServe(addr string, handler http.Handler) error {
	f.Called(addr, handler)
	return nil
}

type mockLogger struct{}

func (f *mockLogger) fatal(v ...interface{}) {}

type mockProvider struct {
	returnFetchErr bool
}

func (f mockProvider) Fetch() (*[]amalgam.Linker, error) {
	return nil, nil
}
func (f mockProvider) Name() string {
	return "Mock Provider"
}
