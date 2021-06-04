package semrel

type RawCommit struct {
	SHA         string
	RawMessage  string
	Annotations map[string]string
}

type Change struct {
	Major       bool
	Minor       bool
	Patch       bool
	Annotations map[string]string
}

type Commit struct {
	SHA         string
	Raw         []string
	Type        string
	Scope       string
	Message     string
	Change      *Change
	Annotations map[string]string
}

type Release struct {
	SHA         string
	Version     string
	Annotations map[string]string
}
