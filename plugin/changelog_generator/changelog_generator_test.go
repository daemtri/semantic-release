package changelog_generator

import (
	"strings"
	"testing"

	"github.com/duanqy/semantic-release/pkg/semrel"
)

func TestDefaultGenerator(t *testing.T) {
	commits := []*semrel.Commit{
		{},
		{SHA: "123456789", Type: "feat", Scope: "app", Message: "commit message"},
		{SHA: "abcd", Type: "fix", Scope: "", Message: "commit message"},
		{SHA: "87654321", Type: "ci", Scope: "", Message: "commit message"},
		{SHA: "43218765", Type: "build", Scope: "", Message: "commit message"},
		{SHA: "12345678", Type: "yolo", Scope: "swag", Message: "commit message"},
		{SHA: "12345678", Type: "chore", Scope: "", Message: "commit message", Raw: []string{"", "BREAKING CHANGE: test"}, Change: &semrel.Change{Major: true}},
		{SHA: "12345679", Type: "chore!", Scope: "user", Message: "another commit message", Raw: []string{"another commit message", "changed ID int into UUID"}, Change: &semrel.Change{Major: true}},
		{SHA: "stop", Type: "chore", Scope: "", Message: "not included"},
	}
	latestRelease := &semrel.Release{SHA: "stop"}
	newVersion := "2.0.0"
	generator := &DefaultChangelogGenerator{}
	changelog := generator.Generate(commits, latestRelease, newVersion)
	if !strings.Contains(changelog, "* **app:** commit message (12345678)") ||
		!strings.Contains(changelog, "* commit message (abcd)") ||
		!strings.Contains(changelog, "#### yolo") ||
		!strings.Contains(changelog, "#### Build") ||
		!strings.Contains(changelog, "#### CI") ||
		!strings.Contains(changelog, "```\nBREAKING CHANGE: test\n```") ||
		strings.Contains(changelog, "not included") {
		t.Fail()
	}
}

func TestEmojiGenerator(t *testing.T) {
	commits := []*semrel.Commit{
		{},
		{SHA: "123456789", Type: "feat", Scope: "app", Message: "commit message"},
		{SHA: "abcd", Type: "fix", Scope: "", Message: "commit message"},
		{SHA: "87654321", Type: "ci", Scope: "", Message: "commit message"},
		{SHA: "43218765", Type: "build", Scope: "", Message: "commit message"},
		{SHA: "12345678", Type: "yolo", Scope: "swag", Message: "commit message"},
		{SHA: "12345678", Type: "chore", Scope: "", Message: "commit message", Raw: []string{"", "BREAKING CHANGE: test"}, Change: &semrel.Change{Major: true}},
		{SHA: "12345679", Type: "chore!", Scope: "user", Message: "another commit message", Raw: []string{"another commit message", "changed ID int into UUID"}, Change: &semrel.Change{Major: true}},
		{SHA: "stop", Type: "chore", Scope: "", Message: "not included"},
	}
	latestRelease := &semrel.Release{SHA: "stop"}
	newVersion := "2.0.0"
	generator := &DefaultChangelogGenerator{emojis: true}
	changelog := generator.Generate(commits, latestRelease, newVersion)
	if !strings.Contains(changelog, "* **app:** commit message (12345678)") ||
		!strings.Contains(changelog, "* commit message (abcd)") ||
		!strings.Contains(changelog, "#### 🎁 Feature") ||
		!strings.Contains(changelog, "#### 🐞 Bug Fixes") ||
		!strings.Contains(changelog, "#### 🔁 CI") ||
		!strings.Contains(changelog, "#### 📦 Build") ||
		!strings.Contains(changelog, "#### 📣 Breaking Changes") ||
		!strings.Contains(changelog, "#### yolo") ||
		!strings.Contains(changelog, "```\nBREAKING CHANGE: test\n```") ||
		strings.Contains(changelog, "not included") {
		t.Fail()
	}
}
