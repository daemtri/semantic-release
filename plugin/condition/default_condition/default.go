package default_condition

import (
	"io/ioutil"
	"strings"

	"github.com/duanqy/semantic-release/pkg/plugin"
)

func init() {
	plugin.RegisterCICondition(&DefaultCI{})
}

func ReadGitHead() string {
	data, err := ioutil.ReadFile(".git/HEAD")
	if err != nil {
		return ""
	}
	return strings.TrimSpace(strings.TrimPrefix(string(data), "ref: refs/heads/"))
}

var CIVERSION = "dev"

type DefaultCI struct {
}

func (d *DefaultCI) Version() string {
	return CIVERSION
}

func (d *DefaultCI) Name() string {
	return "default"
}

func (d *DefaultCI) RunCondition(map[string]string) error {
	return nil
}

func (d *DefaultCI) GetCurrentBranch() string {
	return ReadGitHead()
}

func (d *DefaultCI) GetCurrentSHA() string {
	return ReadGitHead()
}
