package runtime

type defaultScope struct {
	definitions map[string]*scope
}

func (scope defaultScope) Definitions() map[string]*scope {
	return scope.definitions
}

func (scope defaultScope) String() string {
	return ""
}

var DefaultScope = defaultScope{
	definitions: map[string]*scope{},
}
