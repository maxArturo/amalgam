package worker_test

import (
	"testing"

	"github.com/maxArturo/amalgam/internal/worker"
	"github.com/stretchr/testify/assert"
)

func TestUtilAddressResolution(t *testing.T) {
	tests := []struct {
		port     string
		expected string
	}{
		{port: "", expected: ":8080"},
		{port: ":9999", expected: ":9999"},
		{port: "9876", expected: ":9876"},
	}

	u := &util.Util{}

	for _, test := range tests {
		res := u.ResolveAddress(test.port)
		assert.Equal(t, test.expected, res, "should be equal")
	}
}
