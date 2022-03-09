package generator

import (
	"testing"
)

func TestExponential_Generate(t *testing.T) {
	t.Parallel()

	verifyGenerate(t, GenerateExponential)
}

func TestExponential_GenerateConstantCapacity(t *testing.T) {
	t.Parallel()

	verifyGenerateConstantCapacity(t, GenerateExponentialConstantCapacity)
}
