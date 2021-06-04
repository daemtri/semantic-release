package plugin

import "github.com/duanqy/semantic-release/pkg/semrel"

type CommitAnalyzer interface {
	Init(map[string]string) error
	Name() string
	Version() string
	Analyze([]*semrel.RawCommit) []*semrel.Commit
}

var (
	commitAnalyzerSet = map[string]CommitAnalyzer{}
)

func RegisterCommitAnalyzer(ca CommitAnalyzer) {
	commitAnalyzerSet[ca.Name()] = ca
}
