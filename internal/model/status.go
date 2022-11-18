package model

import (
	"time"
)

type StatusType int

const (
	CheckResponseSuccess = 1
	CheckResponseFail    = 2
)

type Metric struct {
	Rdt int64 // response duration time

}

type ResponseTime struct {
	DurationMs int64     // response duration time (ms)
	RTime      time.Time // response time
}
