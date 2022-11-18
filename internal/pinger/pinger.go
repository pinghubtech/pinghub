package pinger

import (
	"context"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/antgubarev/pingbot/internal/model"
	"github.com/antgubarev/pingbot/internal/storage"
	"github.com/sirupsen/logrus"
)

type Pinger struct {
	logger          *logrus.Logger
	targetStorage   storage.TargetStorage
	responseStorage storage.ResponseStorage
}

func NewPinger(logger *logrus.Logger,
	targetStorage storage.TargetStorage,
	responseStorage storage.ResponseStorage) *Pinger {
	return &Pinger{
		logger:          logger,
		targetStorage:   targetStorage,
		responseStorage: responseStorage}
}

func (p *Pinger) Run(ctx context.Context) error {
	ticker := time.NewTicker(time.Second)

	targets, err := p.targetStorage.GetAll(ctx)
	if err != nil {

		return fmt.Errorf("load targets: %v", err)
	}

	wg := sync.WaitGroup{}
	for {
		select {
		case <-ticker.C:
			for _, target := range targets {
				lastResponse, err := p.responseStorage.GetLastResonseTime(ctx, target.Id)
				if err != nil {
					p.logger.Errorf("get last response for targer %s: %v", target.Id, err)
					continue
				}
				if !p.isIterationTimeFinished(lastResponse, *target.IterationTime) {
					continue
				}
				wg.Add(1)
				go func(t model.Target) {
					p.logger.Debugf("ping target %s", t.Id)
					status, err := p.ping(&t)
					if err != nil {
						p.logger.Errorf("ping target: %v", err)
						return
					}
					if err := p.responseStorage.Save(ctx, status, &t); err != nil {
						p.logger.Errorf("save response: %v", err)
						return
					}
				}(*target)
			}

		case <-ctx.Done():
			p.logger.Info("Waiting finishing ping requests...")
			wg.Wait()
			p.logger.Info("All requests are finished")
			return nil
		}
	}
}

func (p *Pinger) isIterationTimeFinished(response *time.Time, iterationTime time.Duration) bool {
	if response == nil {
		return true
	}
	return response.Add(iterationTime).Before(time.Now())
}

func (p *Pinger) ping(target *model.Target) (*model.ResponseTime, error) {
	client := http.Client{
		Timeout: *target.TimeOut,
	}

	startedAt := time.Now()
	resp, err := client.Get(target.URL)
	if err != nil {
		return nil, fmt.Errorf("send request: %v", err)
	}
	defer resp.Body.Close()

	return &model.ResponseTime{
		DurationMs: time.Now().Sub(startedAt).Milliseconds(), //nolint:gosimple
		RTime:      startedAt,
	}, nil
}
