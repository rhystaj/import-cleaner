package types

type OptionalBool int

const (
	OBNil OptionalBool = iota
	OBFalse
	OBTrue
)

func (ob OptionalBool) EqualsBool(boolValue bool) bool {
	return ob == OBTrue && boolValue || ob == OBFalse && !boolValue
}
