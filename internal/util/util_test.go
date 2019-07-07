package util

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type mockOSEnvFetcher struct {
	mockPort string
}

func (f *mockOSEnvFetcher) GetEnvPort() string {
	return f.mockPort
}

func TestUtilAddressResolution(t *testing.T) {
	tests := []struct {
		port     string
		expected string
	}{
		{port: "", expected: ":8080"},
		{port: ":9999", expected: ":9999"},
		{port: "9876", expected: ":9876"},
	}

	u := New()

	for _, test := range tests {
		res := u.ResolveAddress(test.port)
		assert.Equal(t, test.expected, res, "should be equal")
	}
}

func TestUtilResolvesPortWithOSEnv(t *testing.T) {
	u := &Util{
		envFetcher: &mockOSEnvFetcher{
			mockPort: "1984",
		},
	}

	res := u.ResolveAddress("")

	assert.Equal(t, ":1984", res, "should be equal")
}
