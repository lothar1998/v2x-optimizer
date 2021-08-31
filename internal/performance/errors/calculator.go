package errors

import "math"

// Info represents values related to an error computation
// that are obtained based on original value and reference value.
type Info struct {
	Value          int
	ReferenceValue int
	AbsoluteError  int
	RelativeError  float64
}

// Calculate calculate errors between original value and reference value.
func Calculate(referenceValue, value int) *Info {
	diff := int(math.Abs(float64(referenceValue - value)))
	return &Info{
		Value:          value,
		ReferenceValue: referenceValue,
		AbsoluteError:  diff,
		RelativeError:  float64(diff) / float64(referenceValue),
	}
}
