package storage

import (
	"context"
	"fmt"
	"io/ioutil" //nolint:staticcheck

	"github.com/antgubarev/pingbot/internal/model"
	"gopkg.in/yaml.v3"
)

type TargetFileProvider struct {
	path string
}

func NewTargetFileProvider() *TargetFileProvider {
	return &TargetFileProvider{path: "config.yaml"}
}
func (f *TargetFileProvider) GetAll(ctx context.Context) ([]*model.Target, error) {
	yamlFile, err := ioutil.ReadFile(f.path)
	if err != nil {

		return nil, fmt.Errorf("reading targets config from %s: %w", f.path, err)
	}

	var targets model.Targets
	err = yaml.Unmarshal(yamlFile, &targets)
	if err != nil {

		return nil, fmt.Errorf("parsing targets config %w", err)
	}

	return targets.Targets, nil
}
