package tif

type TPrimitiveStringer interface {
	ParseString(string) (any, error)
}

type TRange struct {
	Name        string
	Description string
	Primitive   TPrimitiveStringer
	Postfix     string
	Min         int
	Max         int
}

type TEnum struct {
	Name        string
	Description string
	Primitive   TPrimitiveStringer
	Enumerators []TEnumEnumerator
}

type TEnumEnumerator struct {
	Key         string
	Description string
	Value       any
}

type TSimple struct {
	Name        string
	Description string
	Primitive   TPrimitiveStringer
}
