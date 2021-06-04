package plugin

import (
	"fmt"

	"github.com/duanqy/semantic-release/pkg/semrel"
)

type Hooks interface {
	Init(map[string]string) error
	Name() string
	Version() string
	Success(commits []*semrel.Commit, prevRelease *semrel.Release, newRelease *semrel.Release, changelog string, repoInfo *RepositoryInfo) error
	NoRelease(reason NoReleaseReason, message string) error
}

var (
	hooksSet = map[string]Hooks{}
)

func RegisterHooks(h Hooks) {
	hooksSet[h.Name()] = h
}

type ChainedHooksExecutor struct {
	HooksChain []Hooks
}

func (c *ChainedHooksExecutor) Success(commits []*semrel.Commit, prevRelease *semrel.Release, newRelease *semrel.Release, changelog string, repoInfo *RepositoryInfo) error {
	for _, h := range c.HooksChain {
		err := h.Success(commits, prevRelease, newRelease, changelog, repoInfo)
		if err != nil {
			return err
		}
	}
	return nil
}

func (c *ChainedHooksExecutor) NoRelease(reason NoReleaseReason, message string) error {
	for _, h := range c.HooksChain {
		err := h.NoRelease(reason, message)
		if err != nil {
			return err
		}
	}
	return nil
}

func (c *ChainedHooksExecutor) Init(conf map[string]string) error {
	for _, h := range c.HooksChain {
		err := h.Init(conf)
		if err != nil {
			return err
		}
	}
	return nil
}

func (c *ChainedHooksExecutor) GetNameVersionPairs() []string {
	ret := make([]string, len(c.HooksChain))
	for i, h := range c.HooksChain {
		ret[i] = fmt.Sprintf("%s@%s", h.Name(), h.Version())
	}
	return ret
}
