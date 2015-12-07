package runtime

type defMap map[string]fnScope

// The Scope is the single object of Fn;
// it is used to represent all runtime values.
type fnScope interface {
	// Returns the definitions accessible in this scope.
	Definitions() defMap

	// Returns the current scope with the value defined.
	Define(id string, value fnScope) (fnScope, error)

	// Returns a string representation of the scope.
	String() string

	// Evalutes the scope as a function.
	Call([]fnScope) (fnScope, error)
}

type Scope struct {
	parent      *fnScope
	definitions defMap
}

func (scope Scope) Definitions() defMap {
	return scope.definitions
}

func (scope Scope) String() string {
	if scope.definitions["value"] != nil {
		return scope.definitions["value"].String()
	} else {
		return "scope{}"
	}
}
