package plugin

import "github.com/duanqy/semantic-release/pkg/semrel"

type ChangelogGenerator interface {
	Init(map[string]string) error
	Name() string
	Version() string
	Generate(commits []*semrel.Commit, latestRelease *semrel.Release, newVersion string) string
}

var (
	changelogGenerators = map[string]ChangelogGenerator{}
)

func RegisterChangelogGenerator(cg ChangelogGenerator) {
	changelogGenerators[cg.Name()] = cg
}
