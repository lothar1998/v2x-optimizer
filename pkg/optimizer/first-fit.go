package optimizer

import (
	"errors"
	"github.com/lothar1998/v2x-optimizer/pkg/data"
)

// FirstFit is an optimizer which implements first fit algorithm expanded
// to solve heterogeneous bin with different weight problem.
type FirstFit struct{}

// Optimize runs firs-fit algorithm on the given data.
func (f FirstFit) Optimize(data *data.Data) (*Result, error) {
	v := len(data.R)
	n := len(data.MRB)
	sequence := make([]int, v)
	leftSpace := data.MRB

	for i := 0; i < v; i++ {
		j := 0
		for ; j < n; j++ {
			if data.R[i][j] <= leftSpace[j] {
				sequence[i] = j
				leftSpace[j] -= data.R[i][j]
				break
			}
		}
		if j == n {
			return nil, errors.New("cannot assign V to any RRH")
		}
	}

	return toResult(sequence, n), nil
}
