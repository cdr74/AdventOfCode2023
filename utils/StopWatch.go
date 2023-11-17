package utils

import "time"

type Stopwatch struct {
	startTime   time.Time
	elapsedTime time.Duration
	isRunning   bool
}

func NewStopwatch() *Stopwatch {
	return &Stopwatch{
		startTime:   time.Now(),
		elapsedTime: 0,
		isRunning:   false,
	}
}

func (s *Stopwatch) GetElapsedTime() time.Duration {
	return s.elapsedTime
}

func (s *Stopwatch) Start() {
	s.isRunning = true
	s.startTime = time.Now()
}

func (s *Stopwatch) Stop() {
	s.isRunning = false
	s.elapsedTime = time.Since(s.startTime)
}
