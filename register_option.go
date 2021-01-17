package godi

type RegisterOption func(option *registerOption)

type registerOption struct {
	scope   provideScope
	types   provideType
	tag     string
	filters map[string][]string
}

func newRegisterOption() registerOption {
	return registerOption{
		scope:   Singleton,
		types:   Normal,
		tag:     DefaultTag,
		filters: nil,
	}
}

func Type(types ...provideType) RegisterOption {
	return func(option *registerOption) {
		for _, t := range types {
			option.types = option.types | t
		}
	}
}

func Scope(scope provideScope) RegisterOption {
	return func(option *registerOption) {
		option.scope = scope
	}
}

func Tag(tag string) RegisterOption {
	return func(option *registerOption) {
		option.tag = tag
	}
}

func Filter(dependency string, tags ...string) RegisterOption {
	return func(option *registerOption) {
		if option.filters == nil {
			option.filters = make(map[string][]string)
		}
		option.filters[dependency] = tags
	}
}
