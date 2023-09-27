package gitget

type option struct {
	targetDir string
	force     bool
	branch    string
	tag       string
	commit    string
}

type optionFn func(*option)

func WithTargetDir(dir string) optionFn {
	return func(o *option) {
		o.targetDir = dir
	}
}

func WithForce(force bool) optionFn {
	return func(o *option) {
		o.force = force
	}
}

func WithBranch(branch string) optionFn {
	return func(o *option) {
		o.branch = branch
	}
}

func WithTag(tag string) optionFn {
	return func(o *option) {
		o.tag = tag
	}
}

func WithCommit(commit string) optionFn {
	return func(o *option) {
		o.commit = commit
	}
}

func newOptions(opts ...optionFn) *option {
	// Default options
	o := &option{}
	// Apply options
	for _, opt := range opts {
		opt(o)
	}
	return o
}
