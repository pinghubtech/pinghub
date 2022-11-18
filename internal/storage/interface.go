package storage

import (
	"context"
	"time"

	"github.com/antgubarev/pingbot/internal/model"
)

type AggregateFnType string

const (
	AggregateFnMinute AggregateFnType = "1m"
	AggregateFnHour   AggregateFnType = "1h"
	AggregateFnDay    AggregateFnType = "1d"
)

type AvgResponseTimeByPeriodItem struct {
	DateTime        time.Time
	AvgResponseTime int64
}

type ResponseStorage interface {
	Save(ctx context.Context, response *model.ResponseTime, target *model.Target) error
	GetResponseTimesByPeriod(ctx context.Context, from, to time.Time, targetId string) ([]model.ResponseTime, error)
	GetLastResonseTime(ctx context.Context, targetId string) (*time.Time, error)
}

type TargetStorage interface {
	GetAll(ctx context.Context) ([]*model.Target, error)
}
