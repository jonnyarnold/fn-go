package runtime

// A number is a scope representing a numeric value.
type number struct {
	value string
}

func (num number) Definitions() map[string]*scope {
	return nil
}

func (num number) String() string {
	return num.value
}

func Number(num string) number {
	return number{value: num}
}
