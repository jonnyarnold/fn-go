package runtime

import (
	. "github.com/jonnyarnold/fn-go/compiler/parser"
	. "github.com/jonnyarnold/fn-go/compiler/tokeniser"
	"io/ioutil"
)

// Executes an import!, adding the scope from the file into the current scope.
func execInternalImport(fileName StringExpression, scope fnScope) EvalResult {
	fileScope, err := scopeOfFile(fileName.Value)
	if err != nil {
		return EvalResult{Error: err}
	}

	// We cast to Scope, in order to only get the definitions from the scope (and not its parent).
	for id, value := range fileScope.(Scope).definitions {
		scope, err = scope.Define(id, value)
		if err != nil {
			return EvalResult{Error: err}
		}
	}

	return EvalResult{
		Value: fileScope,
		Scope: scope,
		Error: nil,
	}
}

// Executes an import and returns the scope of the file.
func execVariableImport(fileName StringExpression, scope fnScope) EvalResult {
	fileScope, err := scopeOfFile(fileName.Value)
	if err != nil {
		return EvalResult{Error: err}
	}

	return EvalResult{
		Value: fileScope,
		Scope: scope,
	}
}

// Returns the Scope of the given file.
func scopeOfFile(fileName string) (fnScope, error) {
	file, err := ioutil.ReadFile(fileName)
	if err != nil {
		return nil, err
	}

	tokens := Tokenise(string(file))
	expressions, errors := Parse(tokens)

	if errors != nil {
		return nil, errors
	}

	result := Execute(expressions)

	if result.Error != nil {
		return nil, result.Error
	}

	return result.Scope, nil
}
