package worker

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

type mockOsSleep struct {
	called   bool
	duration time.Duration
}

func (f *mockOsSleep) sleep(d time.Duration) {
	f.called = true
	f.duration = d
}

func TestSleepProvider_sleeps(t *testing.T) {
	mockSleep := &mockOsSleep{}
	mockSource := &source{}
	s := &sleepProvider{
		osSleeper: mockSleep,
	}

	done := make(chan *source)

	go s.sleep(mockSource, done, 1*time.Microsecond)
	resSource := <-done

	assert.True(t, mockSleep.called)
	assert.Equal(t, mockSource, resSource)
}

func TestSleepProvider_sleepsSources(t *testing.T) {
	mockSleep := &mockOsSleep{}
	mockSources := []source{
		source{},
		source{},
		source{},
	}
	s := &sleepProvider{
		osSleeper: mockSleep,
	}

	done := make(chan *source)
	pending := make(chan *source)

	for _, v := range mockSources {
		s.sleepSources(done, pending, 1*time.Microsecond)
		source := &v
		done <- source
		resSource := <-pending
		assert.True(t, mockSleep.called)
		assert.Equal(t, source, resSource)

		mockSleep.called = false
	}
}
