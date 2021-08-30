package errors

import (
	"reflect"
	"testing"
)

func TestCalculate(t *testing.T) {
	type args struct {
		referenceValue int
		value          int
	}
	tests := []struct {
		name string
		args args
		want *Info
	}{
		{
			"should calculate info for referenceValue greater than value",
			args{referenceValue: 12, value: 4},
			&Info{
				Value:          4,
				ReferenceValue: 12,
				AbsoluteError:  8,
				RelativeError:  float64(8) / float64(12),
			},
		},
		{
			"should calculate info for referenceValue lower than value",
			args{referenceValue: 3, value: 7},
			&Info{
				Value:          7,
				ReferenceValue: 3,
				AbsoluteError:  4,
				RelativeError:  float64(4) / float64(3),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Calculate(tt.args.referenceValue, tt.args.value); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Calculate() = %v, want %v", got, tt.want)
			}
		})
	}
}
