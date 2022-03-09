package generator

import (
	"testing"
)

func TestNormal_Generate(t *testing.T) {
	t.Parallel()

	verifyGenerate(t, GenerateNormal)
}

func TestNormal_GenerateConstantCapacity(t *testing.T) {
	t.Parallel()

	verifyGenerateConstantCapacity(t, GenerateNormalConstantCapacity)
}
