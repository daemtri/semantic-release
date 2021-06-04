package plugin

type CICondition interface {
	Name() string
	Version() string
	RunCondition(map[string]string) error
	GetCurrentBranch() string
	GetCurrentSHA() string
}

var (
	ciConditionSet = map[string]CICondition{}
)

func RegisterCICondition(cc CICondition) {
	ciConditionSet[cc.Name()] = cc
}
