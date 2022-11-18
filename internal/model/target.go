package model

import (
	"errors"
	"fmt"
	"net/url"
	"time"
)

type RequestType string

type Provider interface {
	Load() (*Targets, error)
}

const (
	GetRequest RequestType = "GET"
)

type Target struct {
	Id            string         `yaml:"id"`
	Name          string         `yaml:"name,omitempty"`
	URL           string         `yaml:"url"`
	RequestType   *RequestType   `yaml:"requestType"`
	Active        bool           `yaml:"active,omitempty"`
	IterationTime *time.Duration `yaml:"iterationTime"`
	TimeOut       *time.Duration `yaml:"timeOut"`
}

func (t *Target) SetDefaults() {
	if t.IterationTime == nil {
		d := time.Second * 5
		t.IterationTime = &d
	}

	if t.RequestType == nil {
		gt := GetRequest
		t.RequestType = &gt
	}

	if t.TimeOut == nil {
		d := time.Second
		t.TimeOut = &d
	}
}

func (t *Target) Validate() error {
	if len(t.Name) == 0 {
		return errors.New("target must have a name")
	}

	if _, err := url.ParseRequestURI(t.URL); err != nil {
		return fmt.Errorf("target url is invalid: %w", err)
	}

	if t.IterationTime.Seconds() <= 0 {
		return errors.New("iteration time must be 1 second or more")
	}

	return nil
}

type Targets struct {
	Targets []*Target `yaml:"targets"`
}

// func LoadTargetsConfiguration(cfg *config.Config) (*Targets, error) {
// 	provider := NewFileProvider(cfg.ConfigFile)
// 	targets, err := provider.Load()
// 	if err != nil {
//
// 		return nil, err
// 	}
//
// 	return targets, nil
// }
