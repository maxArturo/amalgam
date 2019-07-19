package worker

import "time"

type osSleeper interface {
	sleep(duration time.Duration)
}

type osSleep struct{}

func (s *osSleep) sleep(d time.Duration) {
	time.Sleep(d)
}

type sleepProvider struct {
	osSleeper
}

func newSleepProvider() *sleepProvider {
	return &sleepProvider{
		osSleeper: &osSleep{},
	}
}

func (provider *sleepProvider) sleepSources(done chan *source, pending chan *source, duration time.Duration) {
	go func() {
		for s := range done {
			newSource := s
			go provider.sleep(newSource, pending, duration+time.Duration(newSource.errCount))
		}
	}()
}

func (provider *sleepProvider) sleep(s *source, done chan *source, sleepDuration time.Duration) {
	provider.osSleeper.sleep(sleepDuration)
	done <- s
}
