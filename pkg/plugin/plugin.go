package plugin

type pluginType string

const (
	PluginProvider           pluginType = "provider"
	PluginFilesUpdater       pluginType = "files_updater"
	PluginHooks              pluginType = "hooks"
	PluginChangelogGenerator pluginType = "changelog_generator"
	PluginCondition          pluginType = "ci_condition"
	PluginCommitAnalyzer     pluginType = "commit_analyzer"
)

type NoReleaseReason int32

const (
	NoReleaseReasonCondition NoReleaseReason = 0
	NoReleaseReasonNoChange  NoReleaseReason = 1
)

type RepositoryInfo struct {
	Owner         string
	Repo          string
	DefaultBranch string
	Private       bool
}

type CreateReleaseConfig struct {
	Changelog  string
	NewVersion string
	Prerelease bool
	Branch     string
	SHA        string
}
