package folder

type opts struct {
	ignoredPaths []string
}

type opt func(opts *opts)

func WithIgnore(ignoredPaths ...string) func(opts *opts) {
	return func(opts *opts) {
		opts.ignoredPaths = ignoredPaths
	}
}
