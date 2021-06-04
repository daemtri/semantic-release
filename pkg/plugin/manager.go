package plugin

import (
	"fmt"

	"github.com/duanqy/semantic-release/pkg/config"
)

type Manager struct {
	config    *config.Config
	discovery *Discovery
}

func NewManager(config *config.Config) (*Manager, error) {
	dis, err := NewDiscovery(config)
	if err != nil {
		return nil, err
	}
	return &Manager{
		config:    config,
		discovery: dis,
	}, nil
}

func (m *Manager) GetCICondition() (CICondition, error) {
	cic, ok := ciConditionSet[m.config.CIConditionPlugin]
	if !ok {
		return nil, fmt.Errorf("plugin %s not found", m.config.CIConditionPlugin)
	}

	return cic, nil
}

func (m *Manager) GetProvider() (Provider, error) {
	prov, err := m.discovery.FindPlugin(PluginProvider, m.config.ProviderPlugin)
	if err != nil {
		return nil, err
	}

	return prov.(Provider), nil
}

func (m *Manager) GetCommitAnalyzer() (CommitAnalyzer, error) {
	ca, err := m.discovery.FindPlugin(PluginCommitAnalyzer, m.config.CommitAnalyzerPlugin)
	if err != nil {
		return nil, err
	}

	return ca.(CommitAnalyzer), nil
}

func (m *Manager) GetChangelogGenerator() (ChangelogGenerator, error) {
	cg, err := m.discovery.FindPlugin(PluginChangelogGenerator, m.config.ChangelogGeneratorPlugin)
	if err != nil {
		return nil, err
	}

	return cg.(ChangelogGenerator), nil
}

func (m *Manager) GetChainedUpdater() (*ChainedUpdater, error) {
	updaters := make([]FilesUpdater, 0)
	for _, pl := range m.config.FilesUpdaterPlugins {
		upd, err := m.discovery.FindPlugin(PluginFilesUpdater, pl)
		if err != nil {
			return nil, err
		}

		updaters = append(updaters, upd.(FilesUpdater))
	}

	u := &ChainedUpdater{
		Updaters: updaters,
	}
	return u, nil
}

func (m *Manager) GetChainedHooksExecutor() (*ChainedHooksExecutor, error) {
	hooksChain := make([]Hooks, 0)
	for _, pl := range m.config.HooksPlugins {
		hp, err := m.discovery.FindPlugin(PluginHooks, pl)
		if err != nil {
			return nil, err
		}

		hooksChain = append(hooksChain, hp.(Hooks))
	}

	return &ChainedHooksExecutor{
		HooksChain: hooksChain,
	}, nil
}

func (m *Manager) Stop() {}

func (m *Manager) EnsureAllPlugins() error {
	pluginMap := map[pluginType]string{
		PluginCondition:          m.config.CIConditionPlugin,
		PluginProvider:           m.config.ProviderPlugin,
		PluginCommitAnalyzer:     m.config.CommitAnalyzerPlugin,
		PluginChangelogGenerator: m.config.ChangelogGeneratorPlugin,
	}
	for t, name := range pluginMap {
		_, err := m.discovery.FindPlugin(t, name)
		if err != nil {
			return err
		}
	}

	for _, pl := range m.config.FilesUpdaterPlugins {
		_, err := m.discovery.FindPlugin(PluginFilesUpdater, pl)
		if err != nil {
			return err
		}
	}

	for _, pl := range m.config.HooksPlugins {
		_, err := m.discovery.FindPlugin(PluginHooks, pl)
		if err != nil {
			return err
		}
	}
	return nil
}
