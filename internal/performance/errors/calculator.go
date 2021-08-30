package errors

import "math"

type Info struct {
	Value          int
	ReferenceValue int
	AbsoluteError  int
	RelativeError  float64
}

func Calculate(referenceValue, value int) *Info {
	diff := int(math.Abs(float64(referenceValue - value)))
	return &Info{
		Value:          value,
		ReferenceValue: referenceValue,
		AbsoluteError:  diff,
		RelativeError:  float64(diff) / float64(referenceValue),
	}
}
