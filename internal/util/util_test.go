package util

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type mockOSEnvFetcher struct {
	mockVal string
}

func (f *mockOSEnvFetcher) getEnv(_ string) string {
	return f.mockVal
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
			mockVal: "1984",
		},
	}

	res := u.ResolveAddress("")

	assert.Equal(t, ":1984", res, "should be equal")
}

func TestGetEmptyInt(t *testing.T) {
	u := &Util{
		envFetcher: &mockOSEnvFetcher{
			mockVal: "",
		},
	}

	val, err := u.GetEnvVarInt("bonanza")
	assert.Equal(t, -1, val, "should be -1 to signify err")
	assert.NotNil(t, err)

}

func TestGetBadInt(t *testing.T) {
	u := &Util{
		envFetcher: &mockOSEnvFetcher{
			mockVal: "bonanza",
		},
	}

	val, err := u.GetEnvVarInt("bonanza")
	assert.Equal(t, -1, val, "should be -1 to signify err")
	assert.NotNil(t, err)

}

func TestGetGoodInt(t *testing.T) {
	u := &Util{
		envFetcher: &mockOSEnvFetcher{
			mockVal: "328",
		},
	}

	val, err := u.GetEnvVarInt("bonanza")
	assert.Equal(t, 328, val, "should be an int")
	assert.Nil(t, err)
}

func TestGetPosInt(t *testing.T) {
	u := &Util{
		envFetcher: &mockOSEnvFetcher{
			mockVal: "-328",
		},
	}

	val, err := u.GetEnvVarInt("bonanza")
	assert.Equal(t, -1, val, "should be -1 to signify not a posint")
	assert.NotNil(t, err)
}
