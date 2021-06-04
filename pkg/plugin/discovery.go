package plugin

import (
	"fmt"

	"github.com/duanqy/semantic-release/pkg/config"
)

type Discovery struct {
	config *config.Config
}

func NewDiscovery(config *config.Config) (*Discovery, error) {
	return &Discovery{config}, nil
}

func (d *Discovery) FindPlugin(t pluginType, name string) (interface{}, error) {
	switch t {
	case PluginCommitAnalyzer:
		if p, ok := commitAnalyzerSet[name]; ok {
			return p, nil
		}
	case PluginHooks:
		if p, ok := hooksSet[name]; ok {
			return p, nil
		}
	case PluginCondition:
		if p, ok := ciConditionSet[name]; ok {
			return p, nil
		}
	case PluginFilesUpdater:
		if p, ok := filesUpdaters[name]; ok {
			return p, nil
		}
	case PluginProvider:
		if p, ok := providers[name]; ok {
			return p, nil
		}
	case PluginChangelogGenerator:
		if p, ok := changelogGenerators[name]; ok {
			return p, nil
		}
	}

	return nil, fmt.Errorf("discovery type: %s, plugin: %s not found", t, name)
}
