package runtime

// A fnString is a scope representing a string.
type fnString struct {
	value string
}

func (str fnString) Definitions() map[string]*scope {
	return nil
}

func (str fnString) String() string {
	return str.value
}

func FnString(str string) fnString {
	return fnString{value: str}
}
