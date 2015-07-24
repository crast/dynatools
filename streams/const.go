package streams

import (
	"time"
)

type IteratorType string

const (
	IteratorAtSequence    IteratorType = "AT_SEQUENCE_NUMBER"
	IteratorAfterSequence IteratorType = "AFTER_SEQUENCE_NUMBER"
	IteratorTrimHorizon   IteratorType = "TRIM_HORIZON"
	IteratorLatest        IteratorType = "LATEST"
)

const (
	noAdjust     = 0
	adjustSlower = -2
	adjustFaster = -3
	stopLoop     = -4
)

type adjustConfig struct {
	start      time.Duration
	min        time.Duration
	max        time.Duration
	increaseBy time.Duration
	decreaseBy time.Duration
}

var consumeAdjust = adjustConfig{
	start:      100 * time.Millisecond,
	min:        100 * time.Millisecond,
	max:        30 * time.Second,
	increaseBy: 50 * time.Millisecond,
	decreaseBy: 45 * time.Millisecond,
}

var shardUpdaterAdjust = adjustConfig{
	start:      10 * time.Second,
	min:        10 * time.Second,
	max:        10 * time.Minute,
	increaseBy: 10 * time.Second,
	decreaseBy: 60 * time.Second,
}
