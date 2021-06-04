package plugin

import (
	"github.com/duanqy/semantic-release/pkg/semrel"
)

type Provider interface {
	Init(map[string]string) error
	Name() string
	Version() string
	GetInfo() (*RepositoryInfo, error)
	GetCommits(fromSha, toSha string) ([]*semrel.RawCommit, error)
	GetReleases(re string) ([]*semrel.Release, error)
	CreateRelease(*CreateReleaseConfig) error
}

var (
	providers = map[string]Provider{}
)

func RegisterProvider(p Provider) {
	providers[p.Name()] = p
}
